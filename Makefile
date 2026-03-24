.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/web-example ./src
	docker-compose build app

run:
	- docker-compose up app

build_tests:
	docker-compose -f docker-compose.test.yaml build tests

test:
	- docker-compose -f docker-compose.test.yaml up tests

deploy: build run
