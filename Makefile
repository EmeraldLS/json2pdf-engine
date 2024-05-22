run: build
	@./bin/main 
	
build:
	@go build -o bin/main main.go


test:
	@go test -v

.PHONY:
	test build run