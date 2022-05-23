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

package main

import (
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitextractor/tasks"
	"github.com/mitchellh/mapstructure"
)

var _ core.PluginMeta = (*GitExtractor)(nil)
var _ core.PluginTask = (*GitExtractor)(nil)

type GitExtractor struct{}

func (plugin GitExtractor) Description() string {
	return "extract infos from git repository"
}

// return all available subtasks, framework will run them for you in order
func (plugin GitExtractor) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectGitRepoMeta,
	}
}

// based on task context and user input options, return data that shared among all subtasks
func (plugin GitExtractor) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.GitExtractorOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}
	err = op.Valid()
	if err != nil {
		return nil, err
	}
	return op, nil
}

func (plugin GitExtractor) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/gitextractor"
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry GitExtractor //nolint
