all:
	test lint

test:
	go test ./... -v

lint:
	golangci-lint run ./... --fast --enable-all

run:
	go run main.go

clean:
	rm -rf build

build:
	go build -o bin/main main.go
