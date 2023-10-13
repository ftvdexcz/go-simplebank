postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:12-alpine

createdb: 
	docker exec -it postgres createdb --username=postgres --owner=postgres simplebank

dropdb: 
	docker exec -it postgres dropdb --username=postgres --owner=postgres simplebank

migrateup: 
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simplebank?sslmode=disable" -verbose up

migrateup1: 
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simplebank?sslmode=disable" -verbose up 1

migratedown: 
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simplebank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simplebank?sslmode=disable" -verbose down 1

sqlc:
	docker run --rm -v "%cd%:/src" -w /src sqlc/sqlc generate

test: 
	go test -v -cover ./...

server: 
	go run main.go
 
.PHONY: postgres createdb dropdb sqlc test server migrateup1 migratedown1