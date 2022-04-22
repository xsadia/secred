migrateup:
	migrate -path migrations/ -database "postgresql://postgres:root@localhost:5432/postgres?sslmode=disable" -verbose up

migratedown:
	migrate -path migrations/ -database "postgresql://postgres:root@localhost:5432/postgres?sslmode=disable" -verbose down 1

run:
	go run cmd/main.go

test:
	go test -v ./api ./internal

coverage:
	go test -cover ./api/ ./internal/

.PHONY: migrateup migratedown