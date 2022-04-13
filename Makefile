migrateup:
	migrate -path migrations/ -database "postgresql://postgres:root@localhost:5432/postgres?sslmode=disable" -verbose up

migratedown:
	migrate -path migrations/ -database "postgresql://postgres:root@localhost:5432/postgres?sslmode=disable" -verbose down

run:
	go run cmd/main.go

test:
	go test -v ./api ./internal

.PHONY: migrateup migratedown