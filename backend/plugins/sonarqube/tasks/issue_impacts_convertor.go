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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/codequality"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	sonarqubeModels "github.com/apache/incubator-devlake/plugins/sonarqube/models"
)

var ConvertIssueImpactsMeta = plugin.SubTaskMeta{
	Name:             "convertIssueImpacts",
	EntryPoint:       ConvertIssueImpacts,
	EnabledByDefault: true,
	Description:      "Convert tool layer table sonarqube_issue_impacts into  domain layer table cq_issue_impacts",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}

func ConvertIssueImpacts(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ISSUES_TABLE)
	cursor, err := db.Cursor(
		dal.From("_tool_sonarqube_issue_impacts p"),
		dal.Join("LEFT JOIN _tool_sonarqube_issues i ON (i.connection_id = p.connection_id AND i.issue_key = p.issue_key)"),
		dal.Where("i.connection_id = ? AND i.project_key = ?", data.Options.ConnectionId, data.Options.ProjectKey))
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&sonarqubeModels.SonarqubeIssue{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(sonarqubeModels.SonarqubeIssueImpact{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			impact := inputRow.(*sonarqubeModels.SonarqubeIssueImpact)
			domainIssueImpact := &codequality.CqIssueImpact{
				CqIssueId:       issueIdGen.Generate(data.Options.ConnectionId, impact.IssueKey),
				SoftwareQuality: impact.SoftwareQuality,
				Severity:        impact.Severity,
			}
			return []interface{}{
				domainIssueImpact,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
