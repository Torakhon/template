CURRENT_DIR=$(shell pwd)

run-migration:
	migrate -path migrations -database "postgresql://postgres:20030505@localhost:5432/v1?sslmode=disable"  up
run-down-migration:
		migrate -path migrations -database "postgresql://postgres:20030505@localhost:5432/v1?sslmode=disable"  down
