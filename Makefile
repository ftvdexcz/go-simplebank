postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:12-alpine

createdb: 
	docker exec -it postgres createdb --username=postgres --owner=postgres simplebank

dropdb: 
	docker exec -it postgres dropdb --username=postgres --owner=postgres simplebank

migrateup: 
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simplebank?sslmode=disable" -verbose up

migratedown: 
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simplebank?sslmode=disable" -verbose down

sqlc:
	docker run --rm -v "%cd%:/src" -w /src sqlc/sqlc generate

test: 
	go test -v -cover ./...

.PHONY: postgres createdb dropdb sqlc test