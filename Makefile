build:
	@go build -o bin/main main.go

run: build
	@./bin/main -api_key="<your-api-key>" -json_file="account_statement.json"

test:
	@go test -v

.PHONY:
	test build run