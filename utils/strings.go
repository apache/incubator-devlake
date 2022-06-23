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

// StringsUniq returns a new String Slice contains deduped elements from `source`
func StringsUniq(source []string) []string {
	book := make(map[string]bool, len(source))
	target := make([]string, 0, len(source))
	for _, str := range source {
		if !book[str] {
			book[str] = true
			target = append(target, str)
		}
	}
	return target
}

// StringsContains checks if  `source` String Slice contains `target` string
func StringsContains(slice []string, target string) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}
