.PHONY: build run test clean

build:
	go build -o bin/knowledge-explorer .

run: build
	./bin/knowledge-explorer

test:
	go test ./...

clean:
	rm -rf bin/
