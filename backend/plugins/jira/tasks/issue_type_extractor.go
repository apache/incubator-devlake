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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

var ExtractIssueTypesMeta = plugin.SubTaskMeta{
	Name:             "extractIssueType",
	EntryPoint:       ExtractIssueType,
	EnabledByDefault: true,
	Description:      "extract Jira issueType",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractIssueType(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	connectionId := data.Options.ConnectionId
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	logger.Info("extract IssueType, connection_id=%d, board_id=%d", connectionId, boardId)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: connectionId,
				BoardId:      boardId,
			},
			Table: RAW_ISSUE_TYPE_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			issueType := &models.JiraIssueType{}
			err := errors.Convert(json.Unmarshal(row.Data, issueType))
			if err != nil {
				return nil, err
			}
			//if issueType.Scope != nil {
			//	// FIXME: skip scope for issueTYpe
			//	return nil, nil
			//}
			issueType.ConnectionId = connectionId
			var result []interface{}
			result = append(result, issueType)
			return result, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
