# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# NOTE: you have to use tabs in this file for make. Not spaces.
# https://stackoverflow.com/questions/920413/make-error-missing-separator
# https://tutorialedge.net/golang/makefiles-for-go-developers/

SHA = $(shell git show -s --format=%h)
TAG = $(shell git tag --points-at HEAD)
VERSION = $(TAG)@$(SHA)

build-plugin:
	@sh scripts/compile-plugins.sh

build-worker:
	go build -ldflags "-X 'github.com/apache/incubator-devlake/version.Version=$(VERSION)'" -o bin/lake-worker ./worker/

build-server:
	go build -ldflags "-X 'github.com/apache/incubator-devlake/version.Version=$(VERSION)'" -o bin/lake

build: build-plugin build-server

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

mock:
	rm -rf mocks
	mockery --all --unroll-variadic=false

test: unit-test e2e-test

unit-test: mock build
	set -e; for m in $$(go list ./... | egrep -v 'test|models|e2e'); do echo $$m; go test -timeout 60s -gcflags=all=-l -v $$m; done

e2e-test: build
	PLUGIN_DIR=$(shell readlink -f bin/plugins) go test -timeout 300s -v ./test/...

e2e-plugins:
	export ENV_PATH=$(shell readlink -f .env); set -e; for m in $$(go list ./plugins/... | egrep 'e2e'); do echo $$m; go test -timeout 300s -gcflags=all=-l -v $$m; done

real-e2e-test:
	PLUGIN_DIR=$(shell readlink -f bin/plugins) go test -v ./e2e/...

lint:
	golangci-lint run

clean:
	@rm -rf bin

restart:
	docker-compose down; docker-compose up -d
