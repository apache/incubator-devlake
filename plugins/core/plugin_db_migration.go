package core

import "github.com/merico-dev/lake/migration"

type Migratable interface {
	MigrationScripts() []migration.Script
}
