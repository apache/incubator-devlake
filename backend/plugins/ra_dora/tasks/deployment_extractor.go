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
	"log"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/ra_dora/models"
)

var _ plugin.SubTaskEntryPoint = ExtractDeployments

func init() {
	RegisterSubtaskMeta(&ExtractDeploymentsMeta)
}

// Task metadata
var ExtractDeploymentsMeta = plugin.SubTaskMeta{
	Name:             "extract_deployments",
	EntryPoint:       ExtractDeployments,
	EnabledByDefault: true,
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	DependencyTables: []string{},
	ProductTables:    []string{RAW_DEPLOYMENT_TABLE},
}

func ExtractDeployments(taskCtx plugin.SubTaskContext) errors.Error {
	log.Println("Iniciando plugin de extract.")

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:    taskCtx,
			Params: taskCtx.GetData(),
			Table:  RAW_DEPLOYMENT_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var deployments models.Deployment

			err := errors.Convert(json.Unmarshal(row.Data, &deployments))
			if err != nil {
				return nil, err
			}

			raDeployment := &models.Deployment{
				Metadata: deployments.Metadata,
				Spec:     deployments.Spec,
				Status:   deployments.Status,
			}

			results := make([]interface{}, 0, 2)
			results = append(results, raDeployment)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	log.Println("Extração de deployments concluída com sucesso!")
	return extractor.Execute()
}
