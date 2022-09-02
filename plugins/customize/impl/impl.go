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

package impl

import (
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/customize/api"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var _ core.PluginMeta = (*Customize)(nil)
var _ core.PluginInit = (*Customize)(nil)
var _ core.PluginApi = (*Customize)(nil)

type Customize struct {
	handlers *api.Handlers
}

func (plugin *Customize) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	basicRes := helper.NewDefaultBasicRes(config, logger, db)
	plugin.handlers = api.NewHandlers(basicRes.GetDal())
	return nil
}

func (plugin Customize) Description() string {
	return "To customize table fields"
}

func (plugin Customize) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/customize"
}

func (plugin *Customize) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		":table/fields": {
			"GET":  plugin.handlers.ListFields,
			"POST": plugin.handlers.CreateFields,
		},
		":table/fields/:field": {
			"DELETE": plugin.handlers.DeleteField,
		},
	}
}
