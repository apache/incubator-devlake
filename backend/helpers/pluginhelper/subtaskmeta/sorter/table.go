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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

type TableSorter struct {
	metas []*plugin.SubTaskMeta
}

func NewTableSorter(metas []*plugin.SubTaskMeta) SubTaskMetaSorter {
	return &TableSorter{metas: metas}
}

func (d *TableSorter) Sort() ([]plugin.SubTaskMeta, errors.Error) {
	return tableTopologicalSort(d.metas)
}

func tableTopologicalSort(metas []*plugin.SubTaskMeta) ([]plugin.SubTaskMeta, errors.Error) {
	constructedMetas := constructDependenciesByTable(metas)
	return dependenciesTopologicalSort(constructedMetas)
}

func constructDependenciesByTable(metas []*plugin.SubTaskMeta) []*plugin.SubTaskMeta {
	// construct map by metas and their produced tables, the key is table, and value is metas
	tableMetasMap := make(map[string][]*plugin.SubTaskMeta)
	for _, item := range metas {
		for _, tableItem := range item.ProductTables {
			if value, ok := tableMetasMap[tableItem]; ok {
				tableMetasMap[tableItem] = append(value, item)
			} else {
				tableMetasMap[tableItem] = []*plugin.SubTaskMeta{item}
			}
		}
	}
	// construct meta dependencies by meta.TableDependencies
	// use noDupMap to deduplicate dependencies of meta
	noDupMap := make(map[*plugin.SubTaskMeta]map[*plugin.SubTaskMeta]any)
	for _, metaItem := range metas {
		// convert dependency tables to dependency metas
		dependenciesMap, ok := noDupMap[metaItem]
		if !ok {
			noDupMap[metaItem] = make(map[*plugin.SubTaskMeta]any)
			dependenciesMap = noDupMap[metaItem]
		}
		for _, tableItem := range metaItem.DependencyTables {
			for _, item := range tableMetasMap[tableItem] {
				dependenciesMap[item] = ""
			}
		}
		metaItem.Dependencies = keys(dependenciesMap)
	}
	return metas
}

func keys[T comparable](raw map[T]any) []T {
	list := make([]T, 0)
	for key := range raw {
		list = append(list, key)
	}
	return list
}
