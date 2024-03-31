CURRENT_DIR=$(shell pwd)

run:
	go run cmd/main.go

proto:
	./scripts/genproto.sh  ${CURRENT_DIR}
	ls genproto/*.pb.go | xargs -n1 -IX bash -c "sed -e '/bool/ s/,omitempty//' X > X.tmp && mv X{.tmp,}"

lint: ## Run golangci-lint with printing to stdout
	golangci-lint -c .golangci.yaml run --build-tags "musl" ./...

