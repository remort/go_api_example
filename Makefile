.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/go-api-example ./src
	docker-compose build

run:
	- docker-compose up app front

deploy: build run

build_tests:
	docker-compose -f docker-compose.test.yaml build tests

test:
	- docker-compose -f docker-compose.test.yaml up tests

load_tests:
	wrk -t8 -c500 -d60s -R1000 -s wrk.lua http://0.0.0.0:8000/api/v1/wallet

run_source:
	docker-compose up -d db
	env $$(xargs < config.env.local) go run ./src/
