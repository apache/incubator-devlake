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
	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/dbt/tasks"
)

var (
	_ core.PluginMeta  = (*Dbt)(nil)
	_ core.PluginTask  = (*Dbt)(nil)
	_ core.PluginModel = (*Dbt)(nil)
)

type Dbt struct{}

func (plugin Dbt) Description() string {
	return "Convert data by dbt"
}

func (plugin Dbt) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.GitMeta,
		tasks.DbtConverterMeta,
	}
}

func (plugin Dbt) GetTablesInfo() []core.Tabler {
	return []core.Tabler{}
}

func (plugin Dbt) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.DbtOptions
	err := helper.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	if op.ProjectPath == "" {
		return nil, errors.Default.New("projectPath is required for dbt plugin")
	}

	if op.ProjectTarget == "" {
		op.ProjectTarget = "dev"
	}

	return &tasks.DbtTaskData{
		Options: &op,
	}, nil
}

func (plugin Dbt) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/dbt"
}
