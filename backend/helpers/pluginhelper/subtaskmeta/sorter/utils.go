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
	"github.com/apache/incubator-devlake/core/plugin"
	"sort"
)

type SubtaskMetaList []*plugin.SubTaskMeta

func (s SubtaskMetaList) Len() int {
	return len(SubtaskMetaList{})
}

func (s SubtaskMetaList) Less(i, j int) bool {
	// correct order is collect, extract, enrich, convert
	switch getSubtaskType(s[i].Name) {
	case prefixCollect:
		return true
	case prefixExtract:
		if getSubtaskType(s[j].Name) == prefixCollect {
			return true
		}
	case prefixEnrich:
		typeJ := getSubtaskType(s[j].Name)
		if typeJ == prefixCollect || typeJ == prefixExtract {
			return true
		}
	case prefixConvert:
		return false
	}
	return false
}

func (s SubtaskMetaList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// topologicalSortSameElements
func topologicalSortSameElements(sameElementsDependencyMap map[string][]string) ([]string, error) {
	sortedKeyList := make([]string, 0)
	for {
		if len(sameElementsDependencyMap) == 0 {
			break
		}
		tmpList := make([]string, 0)
		for key, item := range sameElementsDependencyMap {
			if len(item) == 0 {
				tmpList = append(tmpList, key)
			}
		}
		if len(tmpList) == 0 {
			return nil, fmt.Errorf("cyclic dependency detected: %v", sameElementsDependencyMap)
		}
		// remove item in dependencies map
		for key, value := range sameElementsDependencyMap {
			if contains(tmpList, key) {
				delete(sameElementsDependencyMap, key)
			} else {
				sameElementsDependencyMap[key] = removeElements(value, tmpList)
			}
		}
		sort.Strings(tmpList)
		sortedKeyList = append(sortedKeyList, tmpList...)
	}
	return sortedKeyList, nil
}

// topologicalSortDifferentElements
func topologicalSortDifferentElements(differentElementsDependenciesMap map[string][]string) ([]string, error) {
	sortedKeyList := make([]string, 0)
	for {
		if len(differentElementsDependenciesMap) == 0 {
			break
		}
		tmpKeyList := make([]string, 0)
		tmpValueItemList := make([]string, 0)
		for key, item := range differentElementsDependenciesMap {
			if len(item) == 0 || len(item) == 1 {
				tmpKeyList = append(tmpKeyList, key)
				if len(item) == 1 {
					tmpValueItemList = append(tmpValueItemList, item[0])
				}
				delete(differentElementsDependenciesMap, key)
			}
		}
		if len(tmpKeyList) == 0 {
			return nil, fmt.Errorf("cyclic dependency detected: %v", differentElementsDependenciesMap)
		}
		for key, item := range differentElementsDependenciesMap {
			differentElementsDependenciesMap[key] = removeElements(item, tmpValueItemList)
		}
		sort.Strings(tmpKeyList)
		sortedKeyList = append(sortedKeyList, tmpKeyList...)
	}
	return sortedKeyList, nil
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

func convertSubtaskMetaPointToStruct(rawList []*plugin.SubTaskMeta) []plugin.SubTaskMeta {
	list := make([]plugin.SubTaskMeta, len(rawList))
	for index, value := range rawList {
		list[index] = *value
	}
	return list
}
