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

package helper

import (
	"fmt"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/utils"
)

// MakePipelinePlanSubtasks generates subtasks list based on sub-task meta information and entities wanted by user
func MakePipelinePlanSubtasks(subtaskMetas []core.SubTaskMeta, entities []string) ([]string, error) {
	subtasks := make([]string, 0)
	if len(entities) == 0 {
		return subtasks, nil
	}
	wanted := make(map[string]bool, len(entities))
	for _, entity := range entities {
		if !utils.StringsContains(core.DOMAIN_TYPES, entity) {
			return nil, fmt.Errorf("invalid entity(domain type): %s", entity)
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
