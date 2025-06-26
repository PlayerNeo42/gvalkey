MAIN_PATH=./cmd/server/main.go
TEST_TIMEOUT=30s

.PHONY: build
build:
	goreleaser release --snapshot --clean

.PHONY: dev
dev:
	go run $(MAIN_PATH)

.PHONY: test
test:
	go test -timeout $(TEST_TIMEOUT) ./...

.PHONY: test-verbose
test-verbose:
	go test -v -race -timeout $(TEST_TIMEOUT) ./...

.PHONY: fmt
fmt:
	golangci-lint fmt

.PHONY: lint
lint: fmt
	golangci-lint run --fix

.PHONY: check
check: lint test
