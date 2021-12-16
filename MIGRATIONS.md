# Migrations (Database)

## Tutorials

https://dev.to/techschoolguru/how-to-write-run-database-migration-in-golang-5h6g
https://github.com/golang-migrate/migrate

## Install

brew install golang-migrate

## Create a migration

migrate create -ext sql -dir db/migration -seq init_schema

## Run migrations

migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank" -verbose up

## How to reset your DB. (https://github.com/golang-migrate/migrate/issues/282#issuecomment-530743258)

NOTE: Use this when you get an error like this "error: Dirty database version 16. Fix and force version."

1. Delete your DB
2. `migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank" -verbose force 1`
3. `migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank" -verbose down`
4. `migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank" -verbose up`


