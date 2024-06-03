run:
	go run cmd/main.go

race:
	go run -race cmd/main.go

fm:
	gofmt -w .

sc:
	go run scripts/script.go