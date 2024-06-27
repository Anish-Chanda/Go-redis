build:
	@cd ./app && go build -o ../bin/server
run: build
	@./bin/server
test: run
	@go text ./...
