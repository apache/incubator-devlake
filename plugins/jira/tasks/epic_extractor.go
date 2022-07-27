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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var _ core.SubTaskEntryPoint = ExtractIssues

var ExtractEpicsMeta = core.SubTaskMeta{
	Name:             "extractEpics",
	EntryPoint:       ExtractEpics,
	EnabledByDefault: true,
	Description:      "extract Jira epics from all boards",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET, core.DOMAIN_TYPE_CROSS},
}

func ExtractEpics(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDal()
	connectionId := data.Options.ConnectionId
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	logger.Info("extract external epic Issues, connection_id=%d, board_id=%d", connectionId, boardId)
	mappings, err := getTypeMappings(data, db)
	if err != nil {
		return err
	}
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraEpicParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_EPIC_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			return extractIssues(data, mappings, true, row)
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}
