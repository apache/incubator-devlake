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

// SliceRemove remove some items in old slice
func SliceRemove[T ~int | ~string](source []T, toRemoves ...T) []T {
	j := 0
	for _, v := range source {
		needRemove := false
		for _, toRemove := range toRemoves {
			if v == toRemove {
				needRemove = true
				break
			}
		}
		if !needRemove {
			source[j] = v
			j++
		}
	}
	return source[:j]
}
