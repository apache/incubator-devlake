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
	docker-compose up -d grafana

compose-down:
	docker-compose down

commit:
	git cz

install:
	go clean --modcache
	go get

test: unit-test e2e-test models-test

unit-test: build
	go test -v $$(go list ./... | grep -v /test/ | grep -v /models/)

models-test:
	TEST=true go test ./models/test -v

e2e-test: build
	TEST=true PLUGIN_DIR=$(shell readlink -f bin/plugins) go test -v ./test/...

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

init-db:
	mysql -u root --protocol=tcp -padmin -h localhost -e "CREATE DATABASE IF NOT EXISTS lake_test;"
	mysql -u root --protocol=tcp -padmin -h localhost -e "CREATE DATABASE IF NOT EXISTS lake;"
	mysql -u root --protocol=tcp -padmin -h localhost -e "CREATE USER IF NOT EXISTS 'merico'@'localhost' IDENTIFIED BY 'merico';"
	mysql -u root --protocol=tcp -padmin -h localhost -e "GRANT ALL PRIVILEGES ON *.* TO 'merico'@'%';"
	mysql -u root --protocol=tcp -padmin -h localhost -e "USE lake; CREATE TABLE IF NOT EXISTS schema_migrations (version bigint NOT NULL DEFAULT 1,dirty tinyint(1) NOT NULL DEFAULT 0, PRIMARY KEY (version)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	mysql -u root --protocol=tcp -padmin -h localhost -e "USE lake_test; CREATE TABLE IF NOT EXISTS schema_migrations (version bigint NOT NULL DEFAULT 1,dirty tinyint(1) NOT NULL DEFAULT 0, PRIMARY KEY (version)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"

migrate-db:
	migrate -path db/migration -database "mysql://root:admin@tcp(localhost:3306)/$(db)" --verbose up

