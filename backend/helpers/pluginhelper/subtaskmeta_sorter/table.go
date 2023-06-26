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
	"sort"
	"strings"

	"github.com/apache/incubator-devlake/core/plugin"
)

type TableSorter struct {
	metas []*plugin.SubTaskMeta
}

func NewTableSorter(metas []*plugin.SubTaskMeta) SubTaskMetaSorter {
	return &TableSorter{metas: metas}
}

func (d *TableSorter) Sort() ([]plugin.SubTaskMeta, error) {
	return dependencyTableTopologicalSort(d.metas)
}

type SubtaskPrefix string

const (
	prefixCollect SubtaskPrefix = "collect"
	prefixExtract SubtaskPrefix = "extract"
	prefixConvert SubtaskPrefix = "convert"
	prefixEnrich  SubtaskPrefix = "enrich"
)

func genClassNameByMetaName(rawName string) (string, error) {
	if strings.HasPrefix(rawName, string(prefixEnrich)) {
		return string(prefixEnrich), nil
	}
	if len(rawName) > 7 {
		return rawName[7:], nil
	}
	return "", fmt.Errorf("got illeagal raw name = %s", rawName)
}

/*
To use this sorter, developers need to ensure the following prerequisites
 1. The collect, extract, enrich, convert suffix names of the same task need to be the same,
    such as collectAccounts, extractAccounts, convertAccounts,
    otherwise the sorting algorithm cannot effectively match different tasks.
 2. For the same task, developers need to ensure that the execution order is collect, extract, enrich, convert,
    otherwise the order cannot be correctly sorted
 3. Different tasks have no order requirements for operations on the same table,
    which is most important, because the sorting algorithm will only sort tasks
    that depend on the same table according to the subtask name
 4. All subtaskmeta names start with lowercase
*/
func dependencyTableTopologicalSort(metas []*plugin.SubTaskMeta) ([]plugin.SubTaskMeta, error) {
	// 1. construct data
	classNameToSubtaskListMap := make(map[string][]*plugin.SubTaskMeta) // use subtask class name to get meta list
	classNameToTableListMap := make(map[string][]string)                // use class name get meta name list
	subtaskNameToDataMap := make(map[string]*plugin.SubTaskMeta)        // use name to get meta

	for _, metaItem := range metas {
		taskClassName, err := genClassNameByMetaName(metaItem.Name)
		if err != nil {
			return nil, err
		}
		if value, ok := classNameToSubtaskListMap[taskClassName]; ok {
			classNameToSubtaskListMap[taskClassName] = append(value, metaItem)
		} else {
			classNameToSubtaskListMap[taskClassName] = []*plugin.SubTaskMeta{metaItem}
		}
		if value, ok := classNameToTableListMap[taskClassName]; ok {
			// check if subtask in one class has different tables define
			if len(value) != len(metaItem.DependencyTables) {
				return nil, fmt.Errorf("got different table list in class %s", taskClassName)
			}
			// check list item in value and metaItem.DependencyTables, make sure it's equal
			sort.Strings(value)
			sort.Strings(metaItem.DependencyTables)
			for index, valueItem := range value {
				if valueItem != metaItem.DependencyTables[index] {
					return nil, fmt.Errorf("got different table list in class %s", taskClassName)
				}
			}
		} else {
			classNameToTableListMap[taskClassName] = metaItem.DependencyTables
		}
		subtaskNameToDataMap[metaItem.Name] = metaItem
	}

	// 2. sort
	sortedNameList, err := topologicalSortDifferentElements(classNameToTableListMap)
	if err != nil {
		return nil, err
	}

	// 3. gen subtask meta list by sorted data
	sortedSubtaskMetaList := make([]plugin.SubTaskMeta, 0)
	for _, nameItem := range sortedNameList {
		value, ok := classNameToSubtaskListMap[nameItem]
		if !ok {
			return nil, fmt.Errorf("failed get subtask list by class name = %s", nameItem)
		}
		tmpList := make([]plugin.SubTaskMeta, len(value))
		for _, subtaskItem := range value {
			if len(value) >= 1 && len(subtaskItem.Name) > 7 {
				switch SubtaskPrefix(subtaskItem.Name[:7]) {
				case prefixCollect:
					tmpList[0] = *subtaskItem
				case prefixExtract:
					tmpList[1] = *subtaskItem
				case prefixConvert:
					if len(value) == 2 {
						// enrich with convert
						tmpList[1] = *subtaskItem
					} else if len(value) == 3 {
						tmpList[2] = *subtaskItem
					} else {
						return nil, fmt.Errorf("got wrong length of list with extract subtask")
					}
				case prefixEnrich:
					if len(value) == 2 {
						tmpList[0] = *subtaskItem
					} else {
						return nil, fmt.Errorf("got wrong length of list with extract subtask")
					}
				default:
					return nil, fmt.Errorf("got wrong length of subtask %v", subtaskItem)
				}
			}
		}
		sortedSubtaskMetaList = append(sortedSubtaskMetaList, tmpList...)
	}
	return sortedSubtaskMetaList, nil
}
