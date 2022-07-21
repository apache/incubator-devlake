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
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	githubModels "github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
	"reflect"
)

var ConvertMilestonesMeta = core.SubTaskMeta{
	Name:             "convertMilestones",
	EntryPoint:       ConvertMilestones,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_milestones into  domain layer table milestones",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

type MilestoneConverterModel struct {
	common.RawDataOrigin
	githubModels.GithubMilestone
	GithubId int
}

func ConvertMilestones(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId
	connectionId := data.Options.ConnectionId
	db := taskCtx.GetDal()
	clauses := []dal.Clause{
		dal.Select("gi.github_id, gm.*"),
		dal.From("_tool_github_issues gi"),
		dal.Join("JOIN _tool_github_milestones gm ON gm.milestone_id = gi.milestone_id"),
		dal.Where("gm.repo_id = ?", repoId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	boardIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubRepo{})
	domainBoardId := boardIdGen.Generate(connectionId, repoId)
	sprintIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubMilestone{})
	issueIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubIssue{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: connectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_MILESTONE_TABLE,
		},
		InputRowType: reflect.TypeOf(MilestoneConverterModel{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			response := inputRow.(*MilestoneConverterModel)
			domainSprintId := sprintIdGen.Generate(connectionId, response.GithubMilestone.MilestoneId)
			domainIssueId := issueIdGen.Generate(connectionId, response.GithubId)
			sprint := &ticket.Sprint{
				DomainEntity:    domainlayer.DomainEntity{Id: domainSprintId},
				Name:            response.GithubMilestone.Title,
				Url:             response.GithubMilestone.URL,
				Status:          response.GithubMilestone.State,
				StartedDate:     &response.GithubMilestone.CreatedAt, //GitHub doesn't give us a "start date"
				EndedDate:       response.GithubMilestone.ClosedAt,
				CompletedDate:   response.GithubMilestone.ClosedAt,
				OriginalBoardID: domainBoardId,
			}
			boardSprint := &ticket.BoardSprint{
				BoardId:  domainBoardId,
				SprintId: domainSprintId,
			}
			sprintIssue := &ticket.SprintIssue{
				SprintId: domainSprintId,
				IssueId:  domainIssueId,
			}
			return []interface{}{sprint, sprintIssue, boardSprint}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
