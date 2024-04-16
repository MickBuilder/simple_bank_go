new_database:
	docker run -d --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root postgres

existing_database:
	docker start postgres 

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrateuplatest:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedownlatest:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate -f db/sqlc/sqlc.yaml

mock:
	mockgen -package mockdb -destination db/mock/repository.go learning.com/golang_backend/db/sqlc/repository Repository

test:
	go test -v -cover -short ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown migrateuplatest migratedownlatest new_migration sqlc mock test server