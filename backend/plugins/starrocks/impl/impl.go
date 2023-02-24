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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/starrocks/tasks"
)

type StarRocks string

// make sure interface is implemented
var _ plugin.PluginMeta = (*StarRocks)(nil)
var _ plugin.PluginTask = (*StarRocks)(nil)
var _ plugin.PluginModel = (*StarRocks)(nil)

func (s StarRocks) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.ExportDataTaskMeta,
	}
}

func (s StarRocks) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.StarRocksConfig
	err := helper.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	if op.BeHost == "" {
		op.BeHost = op.Host
	}
	return &op, nil
}

func (s StarRocks) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{}
}

func (s StarRocks) Description() string {
	return "Sync data from database to StarRocks"
}

func (s StarRocks) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/starrocks"
}
