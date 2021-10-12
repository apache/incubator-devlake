package services

import (
	"github.com/merico-dev/lake/plugins/core"
)

/**
return value
{
	"jira": {
		"sources": {
			"POST": *ApiResourceHandler
		}
	}
}
*/
func GetPluginsApiResources() (map[string]map[string]map[string]core.ApiResourceHandler, error) {
	res := make(map[string]map[string]map[string]core.ApiResourceHandler)
	for pluginName, pluginEntry := range core.AllPlugins() {
		res[pluginName] = pluginEntry.ApiResources()
	}
	return res, nil
}
