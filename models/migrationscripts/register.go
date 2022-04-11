package migrationscripts

import "github.com/merico-dev/lake/migration"

// RegisterAll register all the migration scripts of framework
func RegisterAll() {
	migration.Register([]migration.Script{new(initSchemas)}, "Framework")
}
