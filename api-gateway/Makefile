CURRENT_DIR=$(shell pwd)

run:
	go run cmd/main.go

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

proto-gen:
	./scripts/gen-proto.sh  ${CURRENT_DIR}
	ls genproto/*.pb.go | xargs -n1 -IX bash -c "sed -e '/bool/ s/,omitempty//' X > X.tmp && mv X{.tmp,}"

lint: ## Run golangci-lint with printing to stdout
	golangci-lint -c .golangci.yaml run --build-tags "musl" ./...
create-migration:
	migrate create -dir migrations -ext sql -seq create_users_table

run-migration:
	migrate -path migrations -database "postgresql://postgres:20030505@localhost:5432/user_db_exam?sslmode=disable"  up
run-down-migration:
		migrate -path migrations -database "postgresql://postgres:20030505@localhost:5432/user_db_exam?sslmode=disable"  down

swag-gen:
	~/go/bin/swag init -g ./api/router.go -o api/docs
