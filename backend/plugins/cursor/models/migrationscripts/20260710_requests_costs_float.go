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

package migrationscripts

import (
	"github.com/apache/devlake/core/context"
	"github.com/apache/devlake/core/errors"
)

type changeCursorRequestsCostsToFloat struct{}

func (script *changeCursorRequestsCostsToFloat) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if !db.HasTable("_tool_cursor_usage_events") {
		return nil
	}
	switch db.Dialect() {
	case "postgres":
		return db.Exec("ALTER TABLE _tool_cursor_usage_events ALTER COLUMN requests_costs TYPE double precision")
	case "mysql":
		return db.Exec("ALTER TABLE _tool_cursor_usage_events MODIFY COLUMN requests_costs DOUBLE")
	default:
		return nil
	}
}

func (*changeCursorRequestsCostsToFloat) Version() uint64 { return 20260710120000 }

func (*changeCursorRequestsCostsToFloat) Name() string {
	return "cursor change requests_costs column to double"
}
