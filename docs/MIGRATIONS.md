# Migrations (Database)

## Summary
Starting in v0.10.0, DevLake provides a lightweight migration tool for executing migration scripts. 
Both framework itself and plugins define their migration scripts in their own migration folder. 
The migration scripts are written with gorm in Golang to support different SQL dialects.


## Migration script
Migration script describes how to do database migration. 
They implement the `Script` interface. 
When DevLake starts, scripts register themselves to the framework by invoking the `Register` function

```go
type Script interface {
	Up(ctx context.Context, db *gorm.DB) error
	Version() uint64
	Name() string
}
```

## Table `migration_history`

The table tracks migration scripts execution and schemas changes. 
From which, DevLake could figure out the current state of database schemas.
## How it Works
1. check `migration_history` table, calculate all the migration scripts need to be executed.
2. sort scripts by Version in ascending order.
3. execute scripts.
4. save results in the `migration_history` table.
