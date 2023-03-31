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

package remote

import (
	"fmt"
	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	pluginCore "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/server/services/remote/models"
	remote "github.com/apache/incubator-devlake/server/services/remote/plugin"
)

var (
	remotePlugins = make(map[string]models.RemotePlugin)
)

func Init(br context.BasicRes) {
	remote.Init(br)
}

func NewRemotePlugin(info *models.PluginInfo) (models.RemotePlugin, errors.Error) {
	if _, ok := remotePlugins[info.Name]; ok {
		return nil, errors.BadInput.New(fmt.Sprintf("plugin %s already registered", info.Name))
	}
	plugin, err := remote.NewRemotePlugin(info)
	if err != nil {
		return nil, err
	}
	forceMigration := config.GetConfig().GetBool("FORCE_MIGRATION")
	err = plugin.RunMigrations(forceMigration)
	if err != nil {
		return nil, err
	}
	err = pluginCore.RegisterPlugin(info.Name, plugin)
	if err != nil {
		return nil, err
	}
	remotePlugins[info.Name] = plugin
	return plugin, nil
}
