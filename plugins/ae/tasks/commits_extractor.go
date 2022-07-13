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

	"github.com/apache/incubator-devlake/plugins/ae/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type ApiCommitsResponse []AeApiCommit

type AeApiCommit struct {
	HexSha      string `json:"hexsha"`
	AnalysisId  string `json:"analysis_id"`
	AuthorEmail string `json:"author_email"`
	DevEq       int    `json:"dev_eq"`
}

func ExtractCommits(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*AeTaskData)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: AeApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_COMMITS_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			apiCommit := &AeApiCommit{}
			err := json.Unmarshal(row.Data, apiCommit)
			if err != nil {
				return nil, err
			}
			return []interface{}{
				&models.AECommit{
					HexSha:      apiCommit.HexSha,
					AnalysisId:  apiCommit.AnalysisId,
					AuthorEmail: apiCommit.AuthorEmail,
					DevEq:       apiCommit.DevEq,
					AEProjectId: data.Options.ProjectId,
				},
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractCommitsMeta = core.SubTaskMeta{
	Name:             "extractCommits",
	EntryPoint:       ExtractCommits,
	EnabledByDefault: true,
	Description:      "Extract raw commit data into tool layer table ae_commits",
}
