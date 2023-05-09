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
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

var _ plugin.SubTaskEntryPoint = ExtractWorklogs

var ExtractWorklogsMeta = plugin.SubTaskMeta{
	Name:             "extractWorklogs",
	EntryPoint:       ExtractWorklogs,
	EnabledByDefault: true,
	Description:      "extract Jira work logs",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractWorklogs(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_WORKLOGS_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var input apiv2models.Input
			err := errors.Convert(json.Unmarshal(row.Input, &input))
			if err != nil {
				return nil, err
			}
			var worklog apiv2models.Worklog
			err = errors.Convert(json.Unmarshal(row.Data, &worklog))
			if err != nil {
				return nil, err
			}
			return []interface{}{worklog.ToToolLayer(data.Options.ConnectionId, &input.UpdateTime)}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
