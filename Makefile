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

SHA ?= $(shell git show -s --format=%h)
TAG ?= $(shell git tag --points-at HEAD)
IMAGE_REPO ?= "apache"
VERSION = $(TAG)@$(SHA)

build-server-image:
	docker build -t $(IMAGE_REPO)/devlake:$(TAG) --build-arg TAG=$(TAG) --build-arg SHA=$(SHA) --file ./Dockerfile .

build-config-ui-image:
	cd config-ui; docker build -t $(IMAGE_REPO)/devlake-config-ui:$(TAG) --file ./Dockerfile .

build-grafana-image:
	cd grafana; docker build -t $(IMAGE_REPO)/devlake-dashboard:$(TAG) --file ./Dockerfile .

build-images: build-server-image build-config-ui-image build-grafana-image

push-server-image: build-server-image
	docker push $(IMAGE_REPO)/devlake:$(TAG)

push-config-ui-image: build-config-ui-image
	docker push $(IMAGE_REPO)/devlake-config-ui:$(TAG)

push-grafana-image: build-grafana-image
        docker push $(IMAGE_REPO)/devlake-dashboard:$(TAG)

push-images: push-server-image push-config-ui-image push-grafana-image

configure:
	docker-compose up config-ui

configure-dev:
	cd config-ui; npm install; npm start;

commit:
	git cz

restart:
	docker-compose down; docker-compose up -d
