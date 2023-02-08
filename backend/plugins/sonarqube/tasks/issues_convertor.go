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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/securitytesting"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	sonarqubeModels "github.com/apache/incubator-devlake/plugins/sonarqube/models"
	"reflect"
)

var ConvertIssuesMeta = plugin.SubTaskMeta{
	Name:             "convertIssues",
	EntryPoint:       ConvertIssues,
	EnabledByDefault: true,
	Description:      "Convert tool layer table sonarqube_issues into  domain layer table st_issues",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_SECURITY_TESTING},
}

func ConvertIssues(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECTS_TABLE)
	cursor, err := db.Cursor(dal.From(sonarqubeModels.SonarqubeIssue{}),
		dal.Where("connection_id = ? and project = ?", data.Options.ConnectionId, data.Options.ProjectKey))
	if err != nil {
		return err
	}
	defer cursor.Close()

	accountIdGen := didgen.NewDomainIdGenerator(&sonarqubeModels.SonarqubeIssue{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(sonarqubeModels.SonarqubeIssue{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			sonarqubeIssue := inputRow.(*sonarqubeModels.SonarqubeIssue)
			domainIssue := &securitytesting.StIssue{
				DomainEntity:      domainlayer.DomainEntity{Id: accountIdGen.Generate(data.Options.ConnectionId, sonarqubeIssue.Key)},
				BatchId:           sonarqubeIssue.BatchId,
				Key:               sonarqubeIssue.Key,
				Rule:              sonarqubeIssue.Rule,
				Severity:          sonarqubeIssue.Severity,
				Component:         sonarqubeIssue.Component,
				Project:           sonarqubeIssue.Project,
				Line:              sonarqubeIssue.Line,
				Status:            sonarqubeIssue.Status,
				Message:           sonarqubeIssue.Message,
				Debt:              sonarqubeIssue.Debt,
				Effort:            sonarqubeIssue.Effort,
				CommitAuthorEmail: sonarqubeIssue.Author,
				Hash:              sonarqubeIssue.Hash,
				Tags:              sonarqubeIssue.Tags,
				Type:              sonarqubeIssue.Type,
				Scope:             sonarqubeIssue.Scope,
				StartLine:         sonarqubeIssue.StartLine,
				EndLine:           sonarqubeIssue.EndLine,
				StartOffset:       sonarqubeIssue.StartOffset,
				EndOffset:         sonarqubeIssue.EndOffset,
				CreationDate:      sonarqubeIssue.CreationDate,
				UpdateDate:        sonarqubeIssue.UpdateDate,
			}
			return []interface{}{
				domainIssue,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
