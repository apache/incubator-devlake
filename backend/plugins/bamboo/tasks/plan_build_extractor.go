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
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
)

var _ plugin.SubTaskEntryPoint = ExtractPlanBuild

func ExtractPlanBuild(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PLAN_BUILD_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,

		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			res := &models.ApiBambooPlanBuild{}
			err := errors.Convert(json.Unmarshal(resData.Data, res))
			if err != nil {
				return nil, err
			}
			plan := &SimplePlan{}
			err = errors.Convert(json.Unmarshal(resData.Input, plan))
			if err != nil {
				return nil, err
			}
			body := res.Convert()
			body.ConnectionId = data.Options.ConnectionId
			body.ProjectKey = data.Options.ProjectKey
			body.PlanKey = plan.PlanKey
			body.Type = data.RegexEnricher.ReturnNameIfMatched(devops.DEPLOYMENT, body.PlanName)
			body.Environment = data.RegexEnricher.ReturnNameIfMatched(devops.PRODUCTION, body.PlanName)

			results := make([]interface{}, 0)
			results = append(results, body)
			// As job build can get more accuracy repo info,
			// we can collect BambooPlanBuildVcsRevision in job_biuld_extractor
			for _, v := range res.VcsRevisions.VcsRevision {
				results = append(results, &models.BambooPlanBuildVcsRevision{
					ConnectionId:   data.Options.ConnectionId,
					PlanBuildKey:   body.PlanBuildKey,
					RepositoryId:   v.RepositoryId,
					RepositoryName: v.RepositoryName,
					VcsRevisionKey: v.VcsRevisionKey,
				})
			}

			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractPlanBuildMeta = plugin.SubTaskMeta{
	Name:             "ExtractPlanBuild",
	EntryPoint:       ExtractPlanBuild,
	EnabledByDefault: true,
	Description:      "Extract raw data into tool layer table bamboo_plan_builds",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}
