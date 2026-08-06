include .env
export

.PHONY: test clean

test:
	go test

clean:
	go clean -cache
