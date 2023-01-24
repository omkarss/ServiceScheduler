test:
	go clean -testcache && go test ./...

start:
	go run cmd/main/main.go

dep:
	go mod download

all: dep start
