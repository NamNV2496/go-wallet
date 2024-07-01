server:
	go run main.go
sqlc:
	sqlc generate
wire:
	wire ./internal/wiring/
test:
	go test ./...

docker-build:
	docker build -t wallet:latest .

# Must run postgres and wallet in a same network
docker-run:
	docker run --name wallet --network wallet-network  -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@172.18.0.3:5432/postgres?sslmode=disable" wallet:latest

docker-inspect:
	docker container inspect go-wallet-$(name)

# example
# docker-inspect-postgres:
# 	docker container inspect go-wallet-postgres-1

.PHONY: server sqlc wire test docker-run docker-build docker-inspect
