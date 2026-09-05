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
	"strings"

	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
	helper "github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/plugins/cursor/models"
)

type memberRecord struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Id        string `json:"id"`
	Role      string `json:"role"`
	IsRemoved bool   `json:"isRemoved"`
}

// ExtractMembers parses raw member records into tool-layer tables.
func ExtractMembers(taskCtx plugin.SubTaskContext) errors.Error {
	data, ok := taskCtx.TaskContext().GetData().(*CursorTaskData)
	if !ok {
		return errors.Default.New("task data is not CursorTaskData")
	}

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Table:   rawMembersTable,
			Options: rawParamsFromTaskData(data),
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			var record memberRecord
			if err := errors.Convert(json.Unmarshal(row.Data, &record)); err != nil {
				return nil, err
			}

			userId := strings.TrimSpace(record.Id)
			if userId == "" {
				userId = strings.TrimSpace(record.Email)
			}
			if userId == "" {
				return nil, nil
			}

			member := &models.CursorMember{
				ConnectionId: data.Options.ConnectionId,
				ScopeId:      data.Options.ScopeId,
				UserId:       userId,
				Email:        strings.TrimSpace(record.Email),
				Name:         strings.TrimSpace(record.Name),
				Role:         strings.TrimSpace(record.Role),
				IsRemoved:    record.IsRemoved,
			}
			return []interface{}{member}, nil
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}
