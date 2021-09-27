package core

import (
	"errors"
	"fmt"
	"strings"
)

// Allowing plugin to know each other

var plugins map[string]Plugin

func RegisterPlugin(name string, plugin Plugin) error {
	if plugins == nil {
		plugins = make(map[string]Plugin)
	}
	plugins[name] = plugin
	return nil
}

func GetPlugin(name string) (Plugin, error) {
	if plugins == nil {
		return nil, errors.New("RegisterPlugin have never been called.")
	}
	if plugin, ok := plugins[name]; ok {
		return plugin, nil
	}
	return nil, fmt.Errorf("Plugin `%s` doesn't exist", name)
}

func AllPlugins() map[string]Plugin {
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
