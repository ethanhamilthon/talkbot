package storage

import (
	cfg "bot/internal/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)


type User struct {
	ID        int64
	Name      string
	PartnerID int64
	Action    string
}

type Storage struct {
	config *cfg.Config
	client *redis.Client
}

var ctx = context.Background()

func New() *Storage {
	config := cfg.New()
	opt, _ := redis.ParseURL(config.GetRedisAddr())
  client := redis.NewClient(opt)
	return &Storage{
		config: config,
		client: client,
	}
}

func (s *Storage) GetUser(UserID int64) (User, error) {
	return s.getUserFromRedis(UserID)
}

func (s *Storage) SetUser(UserID int64, UserName string) error {
	_, err := s.GetUser(UserID)
	if err == nil {
		return errors.New("User already exists")
	}

	user := User{
		ID:        UserID,
		Name:      UserName,
		PartnerID: 0,
		Action:    "waiting",
	}

	return s.setUserToRedis(user)
}

func (s *Storage) CreateChat(UserID int64, PartnerID int64) error {
	user, err := s.getUserFromRedis(UserID)
	if err != nil {
		return err
	}

	partner, err := s.getUserFromRedis(PartnerID)
	if err != nil {
		return err
	}

	user.PartnerID = PartnerID
	user.Action = "chat"
	partner.PartnerID = UserID
	partner.Action = "chat"

	err = s.setUserToRedis(user)
	if err != nil {
		return err
	}

	err = s.setUserToRedis(partner)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetPartner(UserID int64) (User, error) {
	user, err := s.getUserFromRedis(UserID)
	if err != nil {
		return User{}, errors.New("User does not exist")
	}

	keys, err := s.client.Keys(ctx, "user:*").Result()
	if err != nil {
		return User{}, err
	}

	for _, key := range keys {
		val, err := s.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		partner, err := s.unmarshalUser(val)
		if err != nil {
			continue
		}

		if partner.Action == "waiting" && partner.ID != user.ID {
			return partner, nil
		}
	}

	return User{}, errors.New("No partner found")
}

func (s *Storage) CleanPartner(UserID int64) error {
	user, err := s.getUserFromRedis(UserID)
	if err != nil {
		return errors.New("User does not exist")
	}

	user.PartnerID = 0
	user.Action = "waiting"

	return s.setUserToRedis(user)
}

func (s *Storage) DeleteUser(UserID int64) error {
	err := s.client.Del(ctx, s.getUserKey(UserID)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) getUserKey(UserID int64) string {
	return "user:" + fmt.Sprint(UserID)
}

func (s *Storage) marshalUser(user User) (string, error) {
	userJSON, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	return string(userJSON), nil
}

func (s *Storage) unmarshalUser(data string) (User, error) {
	var user User
	err := json.Unmarshal([]byte(data), &user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *Storage) getUserFromRedis(UserID int64) (User, error) {
	val, err := s.client.Get(ctx, s.getUserKey(UserID)).Result()
	if err == redis.Nil {
		return User{}, errors.New("User does not exist")
	} else if err != nil {
		return User{}, err
	}

	return s.unmarshalUser(val)
}

func (s *Storage) setUserToRedis(user User) error {
	userJSON, err := s.marshalUser(user)
	if err != nil {
		return err
	}

	err = s.client.Set(ctx, s.getUserKey(user.ID), userJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
