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

var ConvertProjectsMeta = plugin.SubTaskMeta{
	Name:             "convertProjects",
	EntryPoint:       ConvertProjects,
	EnabledByDefault: true,
	Description:      "Convert tool layer table sonarqube_projects into  domain layer table projects",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_SECURITY_TESTING},
}

func ConvertProjects(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECTS_TABLE)
	cursor, err := db.Cursor(dal.From(sonarqubeModels.SonarqubeProject{}),
		dal.Where("connection_id = ? and `key` = ?", data.Options.ConnectionId, data.Options.ProjectKey))
	if err != nil {
		return err
	}
	defer cursor.Close()

	accountIdGen := didgen.NewDomainIdGenerator(&sonarqubeModels.SonarqubeProject{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(sonarqubeModels.SonarqubeProject{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			sonarqubeProject := inputRow.(*sonarqubeModels.SonarqubeProject)
			domainProject := &securitytesting.StProject{
				DomainEntity:     domainlayer.DomainEntity{Id: accountIdGen.Generate(data.Options.ConnectionId, sonarqubeProject.Key)},
				Key:              sonarqubeProject.Key,
				Name:             sonarqubeProject.Name,
				Qualifier:        sonarqubeProject.Qualifier,
				Visibility:       sonarqubeProject.Visibility,
				LastAnalysisDate: sonarqubeProject.LastAnalysisDate,
				Revision:         sonarqubeProject.Revision,
			}
			return []interface{}{
				domainProject,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
