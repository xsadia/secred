migrateup:
	migrate -path migrations/ -database "postgresql://postgres:root@localhost:5432/postgres?sslmode=disable" -verbose up

migratedown:
	migrate -path migrations/ -database "postgresql://postgres:root@localhost:5432/postgres?sslmode=disable" -verbose down

.PHONY: migrateup migratedown