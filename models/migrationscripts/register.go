package migrationscripts

import "github.com/merico-dev/lake/migration"

// RegisterAll register all the migration scripts of framework
func RegisterAll() {
	migration.Register([]migration.Script{new(initSchemas), new(updateSchemas20220505), new(updateSchemas20220507)}, "Framework")
}
