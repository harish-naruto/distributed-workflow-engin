build:
	go build -o ./bin/server ./server

run: build
	./bin/server