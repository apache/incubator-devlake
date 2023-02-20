/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package runner

import (
	"fmt"
	"io/fs"
	"path/filepath"
	goplugin "plugin"
	"strconv"
	"strings"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/server/services/remote"
)

// LoadPlugins load plugins from local directory
func LoadPlugins(basicRes context.BasicRes) errors.Error {
	remote_plugins_enabled, err := strconv.ParseBool(basicRes.GetConfig("ENABLE_REMOTE_PLUGINS"))
	if err == nil && remote_plugins_enabled {
		remote.Init(basicRes)
	}

	pluginsDir := basicRes.GetConfig("PLUGIN_DIR")
	walkErr := filepath.WalkDir(pluginsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fileName := d.Name()
		if strings.HasSuffix(fileName, ".so") && fileName != ".so" {
			pluginName := fileName[0 : len(d.Name())-3]
			plug, loadErr := goplugin.Open(path)
			if loadErr != nil {
				return loadErr
			}
			symPluginEntry, pluginEntryError := plug.Lookup("PluginEntry")
			if pluginEntryError != nil {
				return pluginEntryError
			}
			pluginMeta, ok := symPluginEntry.(plugin.PluginMeta)
			if !ok {
				return errors.Default.New(fmt.Sprintf("%s PluginEntry must implement PluginMeta interface", pluginName))
			}
			if pluginEntry, ok := symPluginEntry.(plugin.PluginInit); ok {
				err = pluginEntry.Init(basicRes)
				if err != nil {
					return err
				}
			}
			err = plugin.RegisterPlugin(pluginName, pluginMeta)
			if err != nil {
				return nil
			}

			basicRes.GetLogger().Info(`plugin loaded %s`, pluginName)
		}
		return nil
	})
	return errors.Convert(walkErr)
}
