package core

import "github.com/apache/incubator-devlake/migration"

type Migratable interface {
	MigrationScripts() []migration.Script
}
