package plugins

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"plugin"
	"strings"

	. "github.com/merico-dev/lake/plugins/core"
)

// store all plugins
var Plugins map[string]Plugin

// load plugins from local directory
func LoadPlugins(pluginsDir string) error {
	Plugins = make(map[string]Plugin)
	dirs, err := ioutil.ReadDir(pluginsDir)
	if err != nil {
		return err
	}
	for _, subDir := range dirs {
		if !subDir.IsDir() {
			continue
		}
		subDirPath := path.Join(pluginsDir, subDir.Name())
		files, err := ioutil.ReadDir(subDirPath)
		if err != nil {
			return err
		}
		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".so") {
				continue
			}
			so := filepath.Join(subDirPath, file.Name())
			plug, err := plugin.Open(so)
			if err != nil {
				return err
			}
			symPluginEntry, err := plug.Lookup("PluginEntry")
			plugEntry, ok := symPluginEntry.(Plugin)
			if !ok {
				return errors.New(fmt.Sprintf("%v PluginEntry must implement Plugin interface", file.Name()))
			}
			Plugins[subDir.Name()] = plugEntry
			break
		}
	}
	return nil
}

func RunPlugin(name string, options map[string]interface{}, progress chan<- float32) error {
	if Plugins == nil {
		return errors.New("Plugins have to be loaded first, please call LoadPlugins beforehand.")
	}
	plugin, ok := Plugins[name]
	if !ok {
		return errors.New(fmt.Sprintf("Unable to find plugin with name %v", name))
	}
	plugin.Execute(options, progress)
	return nil
}
