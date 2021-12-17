# Migrations (Database)

## Tutorials

https://dev.to/techschoolguru/how-to-write-run-database-migration-in-golang-5h6g
https://github.com/golang-migrate/migrate

## Install

`brew install golang-migrate`

or...

`https://github.com/golang-migrate/migrate/blob/5bf05dc3236ef077e5927c9ca9ca02857a87c582/cmd/migrate/README.md`

## Create a migration

migrate create -ext sql -dir db/migration -seq init_schema

## Run migrations

`migrate -path db/migration -database "mysql://root:admin@tcp(localhost:3306)/lake" -verbose up`

## How to reset your DB. (https://github.com/golang-migrate/migrate/issues/282#issuecomment-530743258)

NOTE: Use this when you get an error like this "error: Dirty database version 16. Fix and force version."

1. Delete your DB
2. `migrate -path db/migration -database "mysql://root:admin@tcp(localhost:3306)/lake" -verbose force 1`
3. `migrate -path db/migration -database "mysql://root:admin@tcp(localhost:3306)/lake" -verbose down`
4. `migrate -path db/migration -database "mysql://root:admin@tcp(localhost:3306)/lake" -verbose up`

## How Migrations Work

1. We are using package "golang-migrate" to run migrations
2. This can be used in CLI, or in the Go program
3. To run a migration, we need two things:
  - Path (directory) to migration scripts
  - Destination (DB connection string) to run migrations on
4. To run these, we need all the migration scripts written.
5. We need "up" scripts, and "down" scripts.
6. Running "up" creates things, "down" removes things.
7. We can also use up and down to alter tables.
8. Currently, we have scripts that create a test DB for us.
9. The DB is created, and all the tables built.
10. We used mysqldump to get a full SQL script for our current DB.
11. ```make models-test``` command does a few things
  - Initializes the test DB using admin permissions (main_test)
  - Runs all tests in models/test folder (insert tests)

