test:
	go test -race ./...

run-example:
	go run ./main.go

run-logstack-example:
	go run ./logstack/main.go

build:
	docker build . -t logserver