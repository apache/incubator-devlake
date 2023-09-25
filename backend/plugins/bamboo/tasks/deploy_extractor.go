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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
)

var _ plugin.SubTaskEntryPoint = ExtractDeploy

func ExtractDeploy(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_DEPLOY_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			res := &models.ApiBambooDeployProject{}
			err := errors.Convert(json.Unmarshal(resData.Data, res))
			if err != nil {
				return nil, err
			}

			var results []interface{}
			if res.PlanKey.Key == data.Options.PlanKey {
				for _, apiEnv := range res.Environments {
					env := new(models.BambooDeployEnvironment)
					env.Convert(&apiEnv)
					env.ConnectionId = data.Options.ConnectionId
					env.PlanKey = res.PlanKey.Key
					results = append(results, env)
				}
			}
			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractDeployMeta = plugin.SubTaskMeta{
	Name:             "ExtractDeploy",
	EntryPoint:       ExtractDeploy,
	EnabledByDefault: true,
	Description:      "Extract raw data into tool layer table _tool_bamboo_deploy_environments",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}
