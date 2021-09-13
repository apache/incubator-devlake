# NOTE: you have to use tabs in this file for make. Not spaces.
# https://stackoverflow.com/questions/920413/make-error-missing-separator
# https://tutorialedge.net/golang/makefiles-for-go-developers/

hello:
	echo "Hello"

build:
	go build

dev:
	@sh ./scripts/dev.sh

run:
	go run main.go

configure:
	cd config-ui; yarn; npm run dev;

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

unit-test:
	@sh ./scripts/unit-test.sh

e2e-test:
	@sh ./scripts/e2e-test.sh

lint:
	golangci-lint run

