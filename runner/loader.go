package runner

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// LoadPlugins load plugins from local directory
func LoadPlugins(pluginsDir string, config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	walkErr := filepath.WalkDir(pluginsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fileName := d.Name()
		if strings.HasSuffix(fileName, ".so") {
			pluginName := fileName[0 : len(d.Name())-3]
			plug, loadErr := plugin.Open(path)
			if loadErr != nil {
				return loadErr
			}
			symPluginEntry, pluginEntryError := plug.Lookup("PluginEntry")
			if pluginEntryError != nil {
				return pluginEntryError
			}
			pluginMeta, ok := symPluginEntry.(core.PluginMeta)
			if !ok {
				return fmt.Errorf("%v PluginEntry must implement PluginMeta interface", pluginName)
			}
			if plugin, ok := symPluginEntry.(core.PluginInit); ok {
				plugin.Init(config, logger, db)
			}
			err = core.RegisterPlugin(pluginName, pluginMeta)
			if err != nil {
				return nil
			}
			logger.Info(`plugin loaded %s`, pluginName)
		}
		return nil
	})
	return walkErr
}
