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
	"sort"
)

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
