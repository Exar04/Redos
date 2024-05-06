run: build
	@./bin/redos --listenAddr :8888

build:
	@go build -o bin/redos .
