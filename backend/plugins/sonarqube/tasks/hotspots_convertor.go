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
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/codequality"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	sonarqubeModels "github.com/apache/incubator-devlake/plugins/sonarqube/models"
)

var ConvertHotspotsMeta = plugin.SubTaskMeta{
	Name:             "convertHotspots",
	EntryPoint:       ConvertHotspots,
	EnabledByDefault: true,
	Description:      "Convert tool layer table sonarqube_hotspots into  domain layer table cq_hotspots",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}

func ConvertHotspots(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_HOTSPOTS_TABLE)
	cursor, err := db.Cursor(dal.From(sonarqubeModels.SonarqubeHotspot{}),
		dal.Where("connection_id = ? and project_key = ?", data.Options.ConnectionId, data.Options.ProjectKey))
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&sonarqubeModels.SonarqubeHotspot{})
	projectIdGen := didgen.NewDomainIdGenerator(&sonarqubeModels.SonarqubeProject{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(sonarqubeModels.SonarqubeHotspot{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			sonarqubeHotspot := inputRow.(*sonarqubeModels.SonarqubeHotspot)
			domainHotspot := &codequality.CqIssue{
				DomainEntity:             domainlayer.DomainEntity{Id: issueIdGen.Generate(data.Options.ConnectionId, sonarqubeHotspot.HotspotKey)},
				Component:                sonarqubeHotspot.Component,
				ProjectKey:               projectIdGen.Generate(data.Options.ConnectionId, sonarqubeHotspot.ProjectKey),
				Line:                     sonarqubeHotspot.Line,
				StartLine:                sonarqubeHotspot.Line,
				Status:                   sonarqubeHotspot.Status,
				Message:                  sonarqubeHotspot.Message,
				CommitAuthorEmail:        sonarqubeHotspot.Author,
				Assignee:                 sonarqubeHotspot.Assignee,
				Rule:                     sonarqubeHotspot.RuleKey,
				CreatedDate:              sonarqubeHotspot.CreationDate,
				UpdatedDate:              sonarqubeHotspot.UpdateDate,
				Type:                     "HOTSPOTS",
				VulnerabilityProbability: sonarqubeHotspot.VulnerabilityProbability,
				SecurityCategory:         sonarqubeHotspot.SecurityCategory,
			}
			return []interface{}{
				domainHotspot,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
