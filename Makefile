# NOTE: you have to use tabs in this file for make. Not spaces.
# https://stackoverflow.com/questions/920413/make-error-missing-separator

hello:
	echo "Hello"

build:
	go build

run:
	go run main.go

compose: 
	docker-compose -f ./devops/docker-compose.yml --project-directory ./ up

compose-down: 
	docker-compose -f ./devops/docker-compose.yml --project-directory ./ down

commit: 
	git cz

install:
	go get