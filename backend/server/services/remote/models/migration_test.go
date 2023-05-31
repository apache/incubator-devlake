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

package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshallMigrationScript(t *testing.T) {
	raw := []byte(`{
		"name": "test",
		"version": 20230420123456,
		"operations": [
			{
				"type": "execute",
				"sql": "SOME SQL",
				"dialect": "mysql"
			},
			{
				"type": "drop_column",
				"table": "t",
				"column": "c"
			},
			{
				"type": "drop_table",
				"table": "t"
			}
		]
	}`)

	var script RemoteMigrationScript
	err := json.Unmarshal(raw, &script)

	assert.NoError(t, err)
	assert.Equal(t, "test", script.name)
	assert.Equal(t, uint64(20230420123456), script.version)
	assert.Len(t, script.operations, 3)

	op1 := script.operations[0].(*ExecuteOperation)
	assert.Equal(t, "SOME SQL", op1.Sql)
	assert.Equal(t, "mysql", *op1.Dialect)

	op2 := script.operations[1].(*DropColumnOperation)
	assert.Equal(t, "t", op2.Table)
	assert.Equal(t, "c", op2.Column)

	op3 := script.operations[2].(*DropTableOperation)
	assert.Equal(t, "t", op3.Table)
}
