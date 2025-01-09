include .env.dev.local

run:
	go run cmd/api/main.go

format:
	go fmt ./...

gen-db-models:
	sqlc generate

migration-down:
	migrate -path infrastructure/postgresdb/sql/migrations -database "$(DATABASE_CONNECTION)" -verbose down

migration-up:
	migrate -path infrastructure/postgresdb/sql/migrations -database "$(DATABASE_CONNECTION)" -verbose up

run-postgres:
	docker-compose -f docker-compose.dev.yml up postgres -d 

stop-postgres:
	docker-compose -f docker-compose.dev.yml stop postgres

run-redis:
	docker-compose -f docker-compose.dev.yml up redis -d 

stop-redis:
	docker-compose -f docker-compose.dev.yml stop redis
