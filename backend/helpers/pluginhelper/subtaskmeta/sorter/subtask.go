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

package sorter

import (
	"fmt"

	"github.com/apache/incubator-devlake/core/errors"

	"github.com/apache/incubator-devlake/core/plugin"
)

type DependencySorter struct {
	metas []*plugin.SubTaskMeta
}

func NewDependencySorter(metas []*plugin.SubTaskMeta) SubTaskMetaSorter {
	return &DependencySorter{metas: metas}
}

func (d *DependencySorter) Sort() ([]plugin.SubTaskMeta, errors.Error) {
	return dependenciesTopologicalSort(d.metas)
}

// stable topological sort
func dependenciesTopologicalSort(metas []*plugin.SubTaskMeta) ([]plugin.SubTaskMeta, errors.Error) {
	// which state will make a cycle
	dependenciesMap := make(map[string][]string)
	nameMetaMap := make(map[string]*plugin.SubTaskMeta)

	for _, item := range metas {
		if item != nil {
			nameMetaMap[item.Name] = item
		}

		if _, ok := dependenciesMap[item.Name]; !ok {
			if len(item.Dependencies) != 0 {
				dependenciesMap[item.Name] = make([]string, 0)
				for _, dependencyItem := range item.Dependencies {
					dependenciesMap[item.Name] = append(dependenciesMap[item.Name], dependencyItem.Name)
				}
			} else {
				dependenciesMap[item.Name] = make([]string, 0)
			}
		} else {
			return nil, errors.Convert(fmt.Errorf("duplicate subtaskmetas detected in list: %s", item.Name))
		}
	}

	// sort
	orderStrList, err := topologicalSortSameElements(dependenciesMap)
	if err != nil {
		return nil, errors.Convert(err)
	}

	// gen list by sorted name list and return
	orderedSubtaskList := make([]plugin.SubTaskMeta, 0)
	for _, item := range orderStrList {
		value, ok := nameMetaMap[item]
		if !ok {
			return nil, errors.Convert(fmt.Errorf("illeagal subtaskmeta detected %s", item))
		}
		orderedSubtaskList = append(orderedSubtaskList, *value)
	}

	return orderedSubtaskList, nil
}
