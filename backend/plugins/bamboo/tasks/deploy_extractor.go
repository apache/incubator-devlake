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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
)

var _ plugin.SubTaskEntryPoint = ExtractDeploy

func ExtractDeploy(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_DEPLOY_TABLE)

	db := taskCtx.GetDal()
	clauses := []dal.Clause{
		dal.Select("plan_key"),
		dal.From(models.BambooPlan{}.TableName()),
		dal.Where("project_key = ? and connection_id=?", data.Options.ProjectKey, data.Options.ConnectionId),
	}
	cursor, err := db.Cursor(
		clauses...,
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	Plans := make(map[string]bool)

	for cursor.Next() {
		Plan := &models.BambooPlan{}
		err = db.Fetch(cursor, Plan)
		if err != nil {
			return err
		}
		Plans[Plan.PlanKey] = true
	}

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			res := &models.ApiBambooDeployProject{}
			err := errors.Convert(json.Unmarshal(resData.Data, res))
			if err != nil {
				return nil, err
			}
			plan := &SimplePlan{}
			err = errors.Convert(json.Unmarshal(resData.Input, plan))

			if err != nil {
				return nil, err
			}

			results := make([]interface{}, 0, len(res.Environments))

			if Plans[res.PlanKey.Key] {
				for _, env := range res.Environments {
					body := &models.BambooDeployEnvironment{}

					body.Convert(&env)
					body.ConnectionId = data.Options.ConnectionId
					body.ProjectKey = data.Options.ProjectKey
					body.PlanKey = res.PlanKey.Key

					results = append(results, body)
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
	Description:      "Extract raw data into tool layer table _tool_bamboo_deploy_environment",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}
