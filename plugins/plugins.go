package plugins

import (
	"errors"
	"fmt"
	"github.com/merico-dev/lake/config"
	"io/fs"
	"path"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/merico-dev/lake/logger"

	. "github.com/merico-dev/lake/plugins/core"
)

// Plugins store all plugins
var Plugins map[string]Plugin

// LoadPlugins load plugins from local directory
func LoadPlugins(pluginsDir string) error {
	Plugins = make(map[string]Plugin)
	walkErr := filepath.WalkDir(pluginsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fileName := d.Name()
		println(fileName, path)
		if strings.HasSuffix(fileName, ".so") {
			pluginName := fileName[0:len(d.Name())-3]
			plug, loadErr := plugin.Open(path)
			if loadErr != nil {
				return loadErr
			}
			symPluginEntry, pluginEntryError := plug.Lookup("PluginEntry")
			if pluginEntryError != nil {
				return pluginEntryError
			}
			plugEntry, ok := symPluginEntry.(Plugin)
			if !ok {
				return fmt.Errorf("%v PluginEntry must implement Plugin interface", pluginName)
			}
			plugEntry.Init()
			logger.Info(`[plugin-core] init a plugin success`, pluginName)
			Plugins[pluginName] = plugEntry
			logger.Info("[plugin-core] plugin loaded", pluginName)
		}
		return nil
	})
	return walkErr
}

func RunPlugin(name string, options map[string]interface{}, progress chan<- float32) error {
	if Plugins == nil {
		return errors.New("plugins have to be loaded first, please call LoadPlugins beforehand")
	}
	p, ok := Plugins[name]
	if !ok {
		return fmt.Errorf("unable to find plugin with name %v", name)
	}
	p.Execute(options, progress)
	return nil
}

func PluginDir() string {
	pluginDir := config.V.GetString("PLUGIN_DIR")
	if !path.IsAbs(pluginDir) {
		wd := config.V.GetString("WORKING_DIRECTORY")
		pluginDir = filepath.Join(wd, pluginDir)
	}
	return pluginDir
}