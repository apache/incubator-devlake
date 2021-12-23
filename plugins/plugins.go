package plugins

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/merico-dev/lake/config"

	"github.com/merico-dev/lake/logger"

	. "github.com/merico-dev/lake/plugins/core"
)

// LoadPlugins load plugins from local directory
func LoadPlugins(pluginsDir string) error {
	walkErr := filepath.WalkDir(pluginsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("KEVIN >>> err", err)
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
			plugEntry, ok := symPluginEntry.(Plugin)
			if !ok {
				return fmt.Errorf("%v PluginEntry must implement Plugin interface", pluginName)
			}
			plugEntry.Init()
			logger.Info(`[plugins] init a plugin success`, pluginName)
			err = RegisterPlugin(pluginName, plugEntry)
			if err != nil {
				return nil
			}
			logger.Info("[plugins] plugin loaded", pluginName)
		}
		return nil
	})
	return walkErr
}

func RunPlugin(name string, options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	plugin, err := GetPlugin(name)
	if err != nil {
		return err
	}
	return plugin.Execute(options, progress, ctx)
}

func PluginDir() string {
	pluginDir := config.V.GetString("PLUGIN_DIR")
	return pluginDir
}
