.PHONY: test run-example run-logstack-example client

test:
	go test -race ./...

run-example:
	go run ./main.go

run-logstack-example:
	go run ./logstack/main.go

build:
	docker build . -t logserver

client: build
	docker run --rm --name logserver -d -v $PWD/logserver.json:/logserver.json -v $PWD/example:/example logserver
	docker cp logserver:/client/dist ./client/
	docker rm -f logserver