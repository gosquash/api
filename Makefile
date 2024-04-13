test:
	@go test ./...

build:
	@go build -o bin/api cmd/api/main.go

run: build
	@go run cmd/api/main.go
