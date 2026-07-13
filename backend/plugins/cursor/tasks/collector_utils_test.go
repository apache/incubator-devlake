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

package tasks

import (
	"encoding/json"
	"testing"

	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/cursor/models"
)

func TestRawParamsFromTaskDataIncludesEndpoint(t *testing.T) {
	taskData := &CursorTaskData{
		Options: &CursorOptions{
			ConnectionId: 1,
			ScopeId:      "team",
		},
		Connection: &models.CursorConnection{
			CursorConn: models.CursorConn{
				RestConnection: helper.RestConnection{Endpoint: models.DefaultEndpoint},
			},
		},
	}

	params := rawParamsFromTaskData(taskData)
	encoded, err := json.Marshal(params.GetParams())
	if err != nil {
		t.Fatalf("marshal params: %v", err)
	}

	expected := `{"ConnectionId":1,"ScopeId":"team","Endpoint":"https://api.cursor.com"}`
	if string(encoded) != expected {
		t.Fatalf("raw params mismatch:\n got: %s\nwant: %s", encoded, expected)
	}
}
