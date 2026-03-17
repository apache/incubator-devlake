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

package services

import (
	"regexp"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
)

// InsertRow FIXME ...
func InsertRow(table string, rows []map[string]interface{}) (int64, errors.Error) {
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(table) {
		return 0, errors.BadInput.New("table name invalid")
	}

	if allowedTables := cfg.GetString("PUSH_API_ALLOWED_TABLES"); allowedTables != "" {
		allow := false
		for _, t := range strings.Split(allowedTables, ",") {
			if strings.TrimSpace(t) == table {
				allow = true
				break
			}
		}
		if !allow {
			return 0, errors.Forbidden.New("table name is not in the allowed list")
		}
	}

	err := db.Create(rows, dal.From(table))
	if err != nil {
		return 0, err
	}
	return 1, nil
}
