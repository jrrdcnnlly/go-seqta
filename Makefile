include .env
export

.PHONY: test cover clean

test:
	go test -cover

cover:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

clean:
	go clean -cache
