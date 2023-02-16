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

	"github.com/apache/incubator-devlake/core/models/domainlayer/codequality"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
)

func ConvertIssueCodeBlocks(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ISSUES_TABLE)

	cursor, err := db.Cursor(
		dal.From("_tool_sonarqube_issue_code_blocks icb"),
		dal.Join("left join _tool_sonarqube_issues i on i.issue_key = icb.issue_key"),
		dal.Where("icb.connection_id = ? and project_key = ?", data.Options.ConnectionId, data.Options.ProjectKey))
	if err != nil {
		return err
	}
	defer cursor.Close()

	idGen := didgen.NewDomainIdGenerator(&models.SonarqubeIssueCodeBlock{})
	issueIdGen := didgen.NewDomainIdGenerator(&models.SonarqubeIssue{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(models.SonarqubeIssueCodeBlock{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			sonarqubeIssueCodeBlock := inputRow.(*models.SonarqubeIssueCodeBlock)
			domainIssueCodeBlock := &codequality.CqIssueCodeBlock{
				DomainEntity: domainlayer.DomainEntity{Id: idGen.Generate(data.Options.ConnectionId, sonarqubeIssueCodeBlock.Id)},
				IssueKey:     issueIdGen.Generate(data.Options.ConnectionId, sonarqubeIssueCodeBlock.IssueKey),
				Component:    sonarqubeIssueCodeBlock.Component,
				Msg:          sonarqubeIssueCodeBlock.Msg,
				StartLine:    sonarqubeIssueCodeBlock.StartLine,
				EndLine:      sonarqubeIssueCodeBlock.EndLine,
				StartOffset:  sonarqubeIssueCodeBlock.StartOffset,
				EndOffset:    sonarqubeIssueCodeBlock.EndOffset,
			}

			return []interface{}{
				domainIssueCodeBlock,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}

var ConvertIssueCodeBlocksMeta = plugin.SubTaskMeta{
	Name:             "convertIssueCodeBlocks",
	EntryPoint:       ConvertIssueCodeBlocks,
	EnabledByDefault: true,
	Description:      "Convert tool layer table sonarqube_issues into domain layer table cq_issue_code_blocks",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}
