APP_NAME=mail-burrow
DOCKER_COMPOSE=docker compose -f docker/docker-compose.yml

.PHONY: run tidy fmt vet test rabbit-up rabbit-down rabbit-logs clean

run:
	go run ./cmd/api

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

rabbit-up:
	$(DOCKER_COMPOSE) up -d

rabbit-down:
	$(DOCKER_COMPOSE) down

rabbit-logs:
	$(DOCKER_COMPOSE) logs -f rabbitmq

clean:
	rm -f database.sqlite3
