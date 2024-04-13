test:
	@go test ./...

build:
	@go build -o bin/api cmd/api/main.go

run: build
	@go run cmd/api/main.go

docker-dev-up:
	docker compose -f ./docker/docker-compose.yml -f ./docker/dev/docker-compose.dev.yml up

docker-dev-down:
	docker compose -f ./docker/docker-compose.yml -f ./docker/dev/docker-compose.dev.yml down
