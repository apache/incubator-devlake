package runner

import (
	"github.com/apache/incubator-devlake/migration"
	"github.com/apache/incubator-devlake/plugins/core"
)

func RegisterMigrationScripts(scripts []migration.Script, comment string, config core.ConfigGetter, logger core.Logger) {
	for _, script := range scripts {
		if s, ok := script.(core.InjectConfigGetter); ok {
			s.SetConfigGetter(config)
		}
		if s, ok := script.(core.InjectLogger); ok {
			s.SetLogger(logger)
		}
	}
	migration.Register(scripts, comment)
}
