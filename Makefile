build:
	@go build -o bin/DFS

run: build
	@./bin/DFS

test:
	@go test -v ./...