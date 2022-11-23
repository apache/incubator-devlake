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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/starrocks/tasks"
)

type StarRocks string

// make sure interface is implemented
var _ core.PluginMeta = (*StarRocks)(nil)
var _ core.PluginTask = (*StarRocks)(nil)
var _ core.PluginModel = (*StarRocks)(nil)

func (s StarRocks) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.LoadDataTaskMeta,
	}
}

func (s StarRocks) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
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

func (s StarRocks) GetTablesInfo() []core.Tabler {
	return []core.Tabler{}
}

func (s StarRocks) Description() string {
	return "Sync data from database to StarRocks"
}

func (s StarRocks) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/starrocks"
}
