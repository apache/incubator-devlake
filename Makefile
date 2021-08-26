# NOTE: you have to use tabs in this file for make. Not spaces.
# https://stackoverflow.com/questions/920413/make-error-missing-separator
# https://tutorialedge.net/golang/makefiles-for-go-developers/

hello:
	echo "Hello"

build:
	go build

dev:
	go build; ./lake

run:
	go run main.go

compose: 
	docker-compose -f ./devops/docker-compose.yml --project-directory ./ up

compose-down: 
	docker-compose -f ./devops/docker-compose.yml --project-directory ./ down

commit: 
	git cz

install:
	go clean --modcache
	go get

test: test-jira unit-test

ci-test:
	go test -v `go list ./... | grep -v /test/`

unit-test: 
	go test -v `go list ./... | grep -v /test/ | grep -v /plugins/`

test-jira:
	go build -buildmode=plugin -o plugins/jira/jira.so plugins/jira/jira.go
	go test ./plugins -v
