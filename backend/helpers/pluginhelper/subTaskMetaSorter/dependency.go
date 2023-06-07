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

package subTaskMetaSorter

import (
	"fmt"
	"github.com/apache/incubator-devlake/core/plugin"
	"sort"
)

type DependencySorter struct {
	metas []*plugin.SubTaskMeta
}

func NewDependencySorter(metas []*plugin.SubTaskMeta) SubTaskMetaSorter {
	return &DependencySorter{metas: metas}
}

func (d *DependencySorter) Sort() ([]plugin.SubTaskMeta, error) {
	return topologicalSort(d.metas)
}

// stable topological sort
func topologicalSort(metas []*plugin.SubTaskMeta) ([]plugin.SubTaskMeta, error) {
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
			return nil, fmt.Errorf("duplicate subtaskmetas detected in list: %s", item.Name)
		}
	}

	orderedSubtaskList := make([]plugin.SubTaskMeta, 0)
	for {
		if len(dependenciesMap) == 0 {
			break
		}

		tmpList := make([]string, 0)
		for key, item := range dependenciesMap {
			if len(item) == 0 {
				tmpList = append(tmpList, key)
			}
		}
		if len(tmpList) == 0 {
			return nil, fmt.Errorf("cyclic dependency detected: %v", dependenciesMap)
		}

		// remove item in dependencies map
		for key, value := range dependenciesMap {
			if contains(tmpList, key) {
				delete(dependenciesMap, key)
			} else {
				dependenciesMap[key] = removeElements(value, tmpList)
			}
		}

		sort.Strings(tmpList)
		// convert item to subtaskmeta by name, and append to orderedSubtaskList
		for _, item := range tmpList {
			value, ok := nameMetaMap[item]
			if !ok {
				return nil, fmt.Errorf("illeagal subtaskmeta detected %s", item)
			}
			orderedSubtaskList = append(orderedSubtaskList, *value)
		}
	}
	return orderedSubtaskList, nil
}

func contains[T comparable](itemList []T, item T) bool {
	for _, newItem := range itemList {
		if item == newItem {
			return true
		}
	}
	return false
}

func removeElements[T comparable](raw, toRemove []T) []T {
	newList := make([]T, 0)
	for _, item := range raw {
		if !contains(toRemove, item) {
			newList = append(newList, item)
		}
	}
	return newList
}
