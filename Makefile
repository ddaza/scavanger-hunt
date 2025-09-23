.PHONY: run build test

run:
	go run .

build:
	mkdir -p bin
	go build -o bin/scavenger-hunt .

test:
	go test ./...

