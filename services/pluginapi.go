package services

import (
	"github.com/merico-dev/lake/plugins/core"
)

/**
return value
{
	"jira": {
		"connections": {
			"POST": *ApiResourceHandler
		}
	}
}
*/
func GetPluginsApiResources() (map[string]map[string]map[string]core.ApiResourceHandler, error) {
	res := make(map[string]map[string]map[string]core.ApiResourceHandler)
	for pluginName, pluginEntry := range core.AllPlugins() {
		if pluginApi, ok := pluginEntry.(core.PluginApi); ok {
			res[pluginName] = pluginApi.ApiResources()
		}
	}
	return res, nil
}
