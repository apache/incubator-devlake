# NOTE: you have to use tabs in this file for make. Not spaces.
# https://stackoverflow.com/questions/920413/make-error-missing-separator
# https://tutorialedge.net/golang/makefiles-for-go-developers/

build-plugin:
	@sh scripts/compile-plugins.sh

build-worker:
	go build -o bin/lake-worker ./worker/

build: build-plugin
	go build -o bin/lake

all: build build-worker

run:
	go run main.go

worker:
	go run worker/*.go

dev: build-plugin run

configure:
	docker-compose up config-ui

configure-dev:
	cd config-ui; npm install; npm start;

commit:
	git cz

test: unit-test e2e-test

unit-test: build
	set -e; for m in $$(go list ./... | egrep -v 'test|models|e2e'); do echo $$m; go test -v $$m; done

e2e-test: build
	PLUGIN_DIR=$(shell readlink -f bin/plugins) go test -v ./test/...

real-e2e-test:
	PLUGIN_DIR=$(shell readlink -f bin/plugins) go test -v ./e2e/...

lint:
	golangci-lint run

clean:
	@rm -rf bin

restart:
	docker-compose down; docker-compose up -d

test-migrateup:
	migrate -path db/migration -database "mysql://merico:merico@localhost:3306/lake" -verbose up

test-migratedown:
	migrate -path db/migration -database "mysql://merico:merico@localhost:3306/lake" -verbose down
