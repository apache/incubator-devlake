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
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiReleasesMeta)
}

var ExtractApiReleasesMeta = plugin.SubTaskMeta{
	Name:             "extractApiReleases",
	EntryPoint:       ExtractApiReleases,
	EnabledByDefault: true,
	Description:      "Extract raw release data into tool layer table",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	DependencyTables: []string{RawReleaseTable},
	ProductTables: []string{
		models.AzuredevopsRelease{}.TableName(),
	},
}

func ExtractApiReleases(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RawReleaseTable)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			results := make([]interface{}, 0, 1)

			releaseApi := &models.AzuredevopsApiRelease{}
			err := errors.Convert(json.Unmarshal(row.Data, releaseApi))
			if err != nil {
				return nil, err
			}

			release := &models.AzuredevopsRelease{
				ConnectionId:          data.Options.ConnectionId,
				AzuredevopsId:         releaseApi.Id,
				ProjectId:             data.Options.ProjectId,
				Name:                  releaseApi.Name,
				Status:                releaseApi.Status,
				ReleaseDefinitionId:   releaseApi.ReleaseDefinition.Id,
				ReleaseDefinitionName: releaseApi.ReleaseDefinition.Name,
				Description:           releaseApi.Description,
				CreatedOn:             releaseApi.CreatedOn,
				ModifiedOn:            releaseApi.ModifiedOn,
			}

			results = append(results, release)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
