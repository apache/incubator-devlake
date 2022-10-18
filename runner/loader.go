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
	"github.com/apache/incubator-devlake/errors"
	"io/fs"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// LoadPlugins load plugins from local directory
func LoadPlugins(pluginsDir string, config *viper.Viper, logger core.Logger, db *gorm.DB) errors.Error {
	walkErr := filepath.WalkDir(pluginsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fileName := d.Name()
		if strings.HasSuffix(fileName, ".so") && fileName != ".so" {
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
				return errors.Default.New(fmt.Sprintf("%s PluginEntry must implement PluginMeta interface", pluginName))
			}
			if plugin, ok := symPluginEntry.(core.PluginInit); ok {
				err = plugin.Init(config, logger, db)
				if err != nil {
					return err
				}
			}
			err = core.RegisterPlugin(pluginName, pluginMeta)
			if err != nil {
				return nil
			}

			logger.Info(`plugin loaded %s`, pluginName)
		}
		return nil
	})
	return errors.Convert(walkErr)
}
