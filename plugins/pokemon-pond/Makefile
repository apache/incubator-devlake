SHELL := /bin/bash

tidy:
	go mod tidy
	go mod vendor

run: tidy
	go build -o pokemon && ./pokemon

