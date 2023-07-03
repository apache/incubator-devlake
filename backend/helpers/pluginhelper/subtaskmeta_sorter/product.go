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

package subtaskmeta_sorter

import (
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"sort"
)

type DependencyAndProductSorter struct {
	metas []*plugin.SubTaskMeta
}

func NewDependencyAndProductSorter(metas []*plugin.SubTaskMeta) SubTaskMetaSorter {
	return &DependencyAndProductSorter{metas: metas}
}

func (d *DependencyAndProductSorter) Sort() ([]plugin.SubTaskMeta, errors.Error) {
	return dependencyAndProductTableTopologicalSort(d.metas)
}

func dependencyAndProductTableTopologicalSort(metas []*plugin.SubTaskMeta) ([]plugin.SubTaskMeta, errors.Error) {
	// 1. construct tmp data struct
	subtaskNameMetaMap := make(map[string]*plugin.SubTaskMeta)
	nameDependencyMap := make(map[string][]string)
	nameProductMap := make(map[string][]string)

	for _, item := range metas {
		subtaskNameMetaMap[item.Name] = item
		nameDependencyMap[item.Name] = item.DependencyTables
		nameProductMap[item.Name] = item.ProductTables
	}

	// 2. topological sort
	sortedNameList := make([]string, 0)
	for {
		if len(sortedNameList) == len(metas) {
			break
		}

		tmpList := make([]string, 0)
		removableDependencyList := make([]string, 0)
		for key, value := range nameDependencyMap {
			if len(value) == 0 {
				tmpList = append(tmpList, key)
				removableDependencyList = append(removableDependencyList, nameProductMap[key]...)
				delete(nameDependencyMap, key)
			}
		}
		if len(removableDependencyList) == 0 {
			return nil, errors.Default.WrapRaw(fmt.Errorf("cyclic dependency detected, "+
				"list[%s] dependencyMap[%s]", sortedNameList, nameDependencyMap))
		}

		for key, value := range nameDependencyMap {
			nameDependencyMap[key] = removeElements(value, removableDependencyList)
		}

		sort.Strings(tmpList)
		sortedNameList = append(sortedNameList, tmpList...)
	}

	// 3. gen subtask meta list by sorted data
	sortedSubtaskMetaList := make([]plugin.SubTaskMeta, 0)
	for _, nameItem := range sortedNameList {
		sortedSubtaskMetaList = append(sortedSubtaskMetaList, *subtaskNameMetaMap[nameItem])
	}
	return sortedSubtaskMetaList, nil
}
