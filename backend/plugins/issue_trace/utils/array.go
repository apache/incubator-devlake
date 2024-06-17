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

package utils

import "strings"

func StringContains(arr []string, item string) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}

	return false
}

func ResolveMultiChangelogs(from, to string) (removedFrom []string, addedTo []string) {
	splitFromItems := strings.Split(from, ",")
	splitToItems := strings.Split(to, ",")
	fromItems := make([]string, 0)
	toItems := make([]string, 0)
	for _, v := range splitFromItems {
		if strings.TrimSpace(v) == "" {
			continue
		}
		fromItems = append(fromItems, strings.TrimSpace(v))
	}
	for _, v := range splitToItems {
		if strings.TrimSpace(v) == "" {
			continue
		}
		toItems = append(toItems, strings.TrimSpace(v))
	}
	for _, v := range fromItems {
		if StringContains(toItems, v) {
			continue
		}
		removedFrom = append(removedFrom, v)
	}
	for _, v := range toItems {
		if StringContains(fromItems, v) {
			continue
		}
		addedTo = append(addedTo, v)
	}
	return
}
