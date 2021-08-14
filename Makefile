build:
	rm -rf bin
	mkdir bin
	go build -o bin

run: build
	clear
	./bin/api-load-test --endpoint "http://localhost:8000" -r 10 -u 2 $(overrides)

dummy_server:
	clear
	cd server && npm run serve