.PHONY: run build test

run:
	go run .

run-dev:
	TWILIO_SKIP_VERIFY=true go run . 

build:
	mkdir -p bin
	go build -o bin/scavenger-hunt .

test:
	go test ./...


tidy: 
	go mod tidy