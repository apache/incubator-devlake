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
	"reflect"
	"strconv"

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	giteaModels "github.com/apache/incubator-devlake/plugins/gitea/models"
)

var ConvertIssuesMeta = core.SubTaskMeta{
	Name:             "convertIssues",
	EntryPoint:       ConvertIssues,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitea_issues into  domain layer table issues",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ConvertIssues(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ISSUE_TABLE)
	repoId := data.Repo.GiteaId

	issue := &giteaModels.GiteaIssue{}
	cursor, err := db.Cursor(
		dal.From(issue),
		dal.Where("repo_id = ? and connection_id=?", repoId, data.Options.ConnectionId),
	)

	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&giteaModels.GiteaIssue{})
	accountIdGen := didgen.NewDomainIdGenerator(&giteaModels.GiteaAccount{})
	boardIdGen := didgen.NewDomainIdGenerator(&giteaModels.GiteaRepo{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(giteaModels.GiteaIssue{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			issue := inputRow.(*giteaModels.GiteaIssue)
			domainIssue := &ticket.Issue{
				DomainEntity:    domainlayer.DomainEntity{Id: issueIdGen.Generate(data.Options.ConnectionId, issue.GiteaId)},
				IssueKey:        strconv.Itoa(issue.Number),
				Title:           issue.Title,
				Description:     issue.Body,
				AssigneeId:      accountIdGen.Generate(data.Options.ConnectionId, issue.AssigneeId),
				AssigneeName:    issue.AssigneeName,
				LeadTimeMinutes: issue.LeadTimeMinutes,
				CreatorId:       accountIdGen.Generate(data.Options.ConnectionId, issue.AuthorId),
				CreatorName:     issue.AuthorName,
				Url:             issue.Url,
				CreatedDate:     &issue.GiteaCreatedAt,
				UpdatedDate:     &issue.GiteaUpdatedAt,
				ResolutionDate:  issue.ClosedAt,
			}
			if issue.State == "closed" {
				domainIssue.Status = ticket.DONE
			} else {
				domainIssue.Status = ticket.TODO
			}
			boardIssue := &ticket.BoardIssue{
				BoardId: boardIdGen.Generate(data.Options.ConnectionId, repoId),
				IssueId: domainIssue.Id,
			}
			return []interface{}{
				domainIssue,
				boardIssue,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
