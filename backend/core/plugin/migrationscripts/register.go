package migrationscripts

import "github.com/apache/incubator-devlake/core/plugin"

// All this is intended for migrations that are common across plugins (generalized plugin migrations)
func All(p plugin.PluginMeta) []plugin.MigrationScript {
	return allScripts(
		newAddRawParamTableForScopes(p),
	)
}

func allScripts(scripts ...plugin.MigrationScript) []plugin.MigrationScript {
	var all []plugin.MigrationScript
	for _, script := range scripts {
		if script != nil {
			all = append(all, script)
		}
	}
	return all
}
