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
TAG ?= $(shell git tag --points-at HEAD)
IMAGE_REPO ?= "apache"
VERSION = $(TAG)@$(SHA)
PYTHON_DIR ?= "./python"

go-dep:
	go install github.com/vektra/mockery/v2@latest
	go install github.com/swaggo/swag/cmd/swag@v1.8.4
	go install github.com/atombender/go-jsonschema/cmd/gojsonschema@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1

python-dep:
	pip install -r python/requirements.txt

dep: go-dep python-dep

swag:
	if [ -z $(PLUGIN) ]; then \
		swag init --parseDependency --parseInternal -o ./server/api/docs -g ./server/api/api.go -g ./plugins/*/api/*.go; \
	elif [ $(PLUGIN) = "none" ]; then \
		swag init --parseDependency --parseInternal -o ./server/api/docs -g ./server/api/api.go;\
	else \
	  	plugins="";\
		for p in $$(echo $(PLUGIN) | tr "," "\n"); do \
			plugins="$$plugins -g ./plugins/$$p/api/*.go"; \
		done;\
		swag init --parseDependency --parseInternal -o ./server/api/docs -g ./server/api/api.go "$$plugins"; \
	fi;\
	echo "visit the swagger document on http://localhost:8080/swagger/index.html";

build-plugin:
	if [ "$(PLUGIN)" = "none" ]; then \
  		echo "Building plugins will be skipped"; \
	elif [ "$(DEBUG)" = "true" ]; then \
  		PLUGIN=$(PLUGIN) sh scripts/compile-plugins.sh -gcflags='all=-N -l'; \
  	else \
		PLUGIN=$(PLUGIN) sh scripts/compile-plugins.sh; \
	fi

build-worker:
	if [ "$(DEBUG)" = "true" ]; then \
  		go build -gcflags='all=-N -l' -ldflags "-X 'github.com/apache/incubator-devlake/core/version.Version=$(VERSION)'" -o bin/lake-worker ./worker/; \
  	else \
		go build -ldflags "-X 'github.com/apache/incubator-devlake/core/version.Version=$(VERSION)'" -o bin/lake-worker ./worker/; \
	fi

build-server: swag
	if [ "$(DEBUG)" = "true" ]; then \
  		go build -gcflags='all=-N -l' -ldflags "-X 'github.com/apache/incubator-devlake/core/version.Version=$(VERSION)'" -o bin/lake ./server/; \
  	else \
		go build -ldflags "-X 'github.com/apache/incubator-devlake/core/version.Version=$(VERSION)'" -o bin/lake ./server/; \
	fi

build-python: #don't mix this with the other build commands
	find ./python/ -name "*.sh" | xargs chmod +x &&\
	sh python/build.sh
	sh python/build.sh python/test

build: build-plugin build-server

all: build build-worker

tap-models:
	chmod +x ./scripts/singer-model-generator.sh
	@sh scripts/singer-model-generator.sh resources/tap/pagerduty.json plugins/pagerduty --all

run:
	go run server/main.go

worker:
	go run worker/*.go

dev: build-plugin build-python run

godev:
	DISABLED_REMOTE_PLUGINS=true make build-plugin run

debug: build-plugin-debug
	dlv debug server/main.go

mock:
	rm -rf mocks
	mockery --recursive --keeptree --dir=./core --output=./mocks/core --unroll-variadic=false --name='.*'
	mockery --recursive --keeptree --dir=./helpers --output=./mocks/helpers --unroll-variadic=false --name='.*'

test: unit-test e2e-test

unit-test: mock unit-test-only python-unit-test

unit-test-only:
	set -e;\
	for m in $$(go list ./... | egrep -v 'test|models|e2e'); do \
  		echo $$m; \
			if ! go test -timeout 60s -v $$m; then \
				exit $$?; \
			fi \
	done; \

build-pydevlake:
	poetry install -C python/pydevlake

python-unit-test: build-pydevlake
	sh python/build.sh python/test &&\
	sh ./python/run_tests.sh

e2e-plugins-test:
	export ENV_PATH=$(shell readlink -f .env);\
	set -e;\
	for m in $$(go list ./plugins/... | egrep 'e2e'); do \
		echo $$m; \
		if !  go test -timeout 300s -gcflags=all=-l -v $$m; then \
		  exit $$?; \
		fi \
	done; \

e2e-test-init:
	export ENV_PATH=$(shell readlink -f .env);\
	set -e;\
	go run ./test/init.go || exit $$?;\

e2e-test-run:
	for m in $$(go list ./test/e2e/... | grep -v manual); do \
		echo $$m; \
		if ! go test -p 1 -timeout 300s -v $$m ; then \
		  exit $$?; \
		fi; \
	done; \

e2e-test: e2e-test-init e2e-test-run

integration-test:
	export ENV_PATH=$(shell readlink -f .env);\
	set -e;\
	go run ./test/init.go || exit $$?;\
	exit_code=0;\
	for m in $$(go list ./test/integration/...); do \
		echo $$m; go test -p 1 -timeout 300s -v $$m || exit_code=$$?; \
    done; \
	exit $$exit_code

lint:
	golangci-lint run

fmt:
	find . -name \*.go | xargs gofmt -s -w -l

clean:
	@rm -rf bin

build-server-image:
	docker build -t $(IMAGE_REPO)/devlake:$(TAG) --build-arg TAG=$(TAG) --build-arg SHA=$(SHA) --file ./Dockerfile .

migration-script-lint:
	go run core/migration/linter/main.go -- $$(find . -path '**/migrationscripts/**.go')
