.PHONY: test run-example run-logstack-example client

test:
	go test -race ./...

run-example:
	go run ./main.go -debug -config ./example/logserver.json

run-example-dynamic:
	go run ./main.go -debug -dynamic -config ./example/logserver.json

build:
	docker build . -t logserver

client: build
	docker run --rm --name logserver -d -v $(PWD)/logserver.json:/logserver.json -v $(PWD)/example:/example logserver
	docker cp logserver:/client/dist ./client/
	docker rm -f logserver

build-fast:
	docker build ./ -f Dockerfile.fast -t logserver

run-nginx:
	docker run \
		--name nginx \
		--rm \
		-v $(PWD)/example/nginx:/web:ro \
		-v $(PWD)/example/nginx/nginx.conf:/etc/nginx/nginx.conf:ro \
		-p 8001:80 \
		-p 8002:81 \
		-u root \
		--entrypoint nginx \
		nginx:1.12-alpine -g 'daemon off;'
