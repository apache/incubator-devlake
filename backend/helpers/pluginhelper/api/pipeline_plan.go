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

package api

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
)

// MakePipelinePlanSubtasks generates subtasks list based on sub-task meta information and entities wanted by user
func MakePipelinePlanSubtasks(subtaskMetas []plugin.SubTaskMeta, entities []string) ([]string, errors.Error) {
	subtasks := make([]string, 0)
	// if no entities specified, use all entities enabled by default
	if len(entities) == 0 {
		entities = plugin.DOMAIN_TYPES
	}
	wanted := make(map[string]bool, len(entities))
	for _, entity := range entities {
		if !utils.StringsContains(plugin.DOMAIN_TYPES, entity) {
			return nil, errors.Default.New(fmt.Sprintf("invalid entity(domain type): %s", entity))
		}
		wanted[entity] = true
	}
	for _, subtaskMeta := range subtaskMetas {
		if !subtaskMeta.EnabledByDefault {
			continue
		}
		for _, neededBy := range subtaskMeta.DomainTypes {
			if wanted[neededBy] {
				subtasks = append(subtasks, subtaskMeta.Name)
				break
			}
		}
	}
	return subtasks, nil
}

func MakePipelinePlanTask(
	pluginName string,
	subtaskMetas []plugin.SubTaskMeta,
	entities []string,
	options interface{},
) (*models.PipelineTask, errors.Error) {
	// get subtasks enabled by default
	cfg := config.GetConfig()
	enableSubtasksByDefault := cfg.GetString("ENABLE_SUBTASKS_BY_DEFAULT")
	enableSubtasksList := strings.Split(enableSubtasksByDefault, ",")
	for s := range subtaskMetas {
		compareName := pluginName + ":" + subtaskMetas[s].Name
		for _, enableSubtask := range enableSubtasksList {
			subtaskInfo := strings.Split(enableSubtask, ":")
			if len(subtaskInfo) > 2 {
				subtaskInfoName := subtaskInfo[0] + ":" + subtaskInfo[1]
				if subtaskInfoName == compareName {
					v, err := strconv.ParseBool(subtaskInfo[2])
					if err != nil {
						break
					}
					subtaskMetas[s].EnabledByDefault = v
				}
			}
		}
	}

	subtasks, err := MakePipelinePlanSubtasks(subtaskMetas, entities)
	if err != nil {
		return nil, err
	}
	op, err := encodeTaskOptions(options)
	if err != nil {
		return nil, err
	}
	return &models.PipelineTask{
		Plugin:   pluginName,
		Subtasks: subtasks,
		Options:  op,
	}, nil
}

func encodeTaskOptions(op interface{}) (map[string]interface{}, errors.Error) {
	var result map[string]interface{}
	err := Decode(op, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}
