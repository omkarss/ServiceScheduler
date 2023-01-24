
test:
	go test ./...

start:
	go run cmd/main/main.go

dep:
	go mod download

lint:
	golangci-lint run --enable-all