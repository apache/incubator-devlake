package core

import (
	"errors"
	"fmt"
	"strings"
)

// Allowing plugin to know each other

var plugins map[string]PluginMeta

func RegisterPlugin(name string, plugin PluginMeta) error {
	if plugins == nil {
		plugins = make(map[string]PluginMeta)
	}
	plugins[name] = plugin
	return nil
}

func GetPlugin(name string) (PluginMeta, error) {
	if plugins == nil {
		return nil, errors.New("RegisterPlugin have never been called.")
	}
	if plugin, ok := plugins[name]; ok {
		return plugin, nil
	}
	return nil, fmt.Errorf("Plugin `%s` doesn't exist", name)
}

func AllPlugins() map[string]PluginMeta {
	return plugins
}

func FindPluginNameBySubPkgPath(subPkgPath string) (string, error) {
	for name, plugin := range plugins {
		if strings.HasPrefix(subPkgPath, plugin.RootPkgPath()) {
			return name, nil
		}
	}
	return "", fmt.Errorf("Unable to find plugin for subPkgPath %s", subPkgPath)
}
