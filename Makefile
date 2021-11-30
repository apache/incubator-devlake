# NOTE: you have to use tabs in this file for make. Not spaces.
# https://stackoverflow.com/questions/920413/make-error-missing-separator
# https://tutorialedge.net/golang/makefiles-for-go-developers/

hello:
	echo "Hello"

build-plugin:
	@sh scripts/compile-plugins.sh

build: build-plugin
	go build -o bin/lake

dev: build
	bin/lake

run:
	go run main.go

configure:
	docker-compose up config-ui

configure-dev:
	cd config-ui; npm install; npm start;

compose:
	docker-compose up grafana

compose-down:
	docker-compose down

commit:
	git cz

install:
	go clean --modcache
	go get

test: unit-test e2e-test

unit-test: build
	go test -v $$(go list ./... | grep -v /test/)

e2e-test: build
	PLUGIN_DIR=$(shell readlink -f bin/plugins) go test -v ./test/...

lint:
	golangci-lint run

clean:
	@rm -rf bin

restart:
	docker-compose down; docker-compose up -d
