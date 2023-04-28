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

package dalgorm

import (
	"gorm.io/gorm/schema"
	"time"
)

// ToDatabaseMap convert the map to a format that can be inserted into a SQL database
func ToDatabaseMap(tableName string, m map[string]any) map[string]any {
	strategy := schema.NamingStrategy{}
	newMap := map[string]any{}
	for k, v := range m {
		k = strategy.ColumnName(tableName, k)
		if str, ok := v.(string); ok {
			t, err := time.Parse("0001-01-01T00:00:00Z", str) //time needs special handling
			if err == nil {
				if t.Second() == 0 { // 0 time is not permitted - just skip the column
					continue
				} else {
					v = t
				}
			}
		}
		newMap[k] = v
	}
	return newMap
}
