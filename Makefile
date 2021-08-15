build:
	rm -rf bin
	mkdir bin
	go build -o bin

test:
	go test -v ./...