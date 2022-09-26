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
	"database/sql"
	"reflect"
	"regexp"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var EnrichTaskEnvMeta = core.SubTaskMeta{
	Name:             "EnrichTaskEnv",
	EntryPoint:       EnrichTasksEnv,
	EnabledByDefault: true,
	Description:      "calculate deployment frequency",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func EnrichTasksEnv(taskCtx core.SubTaskContext) (err errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*DoraTaskData)
	repoId := data.Options.RepoId

	productionNamePattern := data.Options.ProductionPattern
	// TODO: STAGE 2
	// stagingNamePattern := data.Options.StagingPattern
	// testingNamePattern := data.Options.TestingPattern
	dataSource := data.Options.DataSoure

	productionNameRegexp, errRegexp := regexp.Compile(productionNamePattern)
	if errRegexp != nil {
		return errors.Default.Wrap(errRegexp, "Regexp compile productionPattern failed")
	}
	// TODO: STAGE 2
	// stagingNameRegexp, errRegexp := regexp.Compile(stagingNamePattern)
	// if errRegexp != nil {
	// 	return errors.Default.Wrap(errRegexp, "Regexp compile stagingPattern failed")
	// }
	// testingNameRegexp, errRegexp := regexp.Compile(testingNamePattern)
	// if errRegexp != nil {
	// 	return errors.Default.Wrap(errRegexp, "Regexp compile testingPattern failed")
	// }

	var cursor *sql.Rows
	if len(dataSource) == 0 {
		cursor, err = db.Cursor(
			dal.From(&devops.CICDTask{}),
			dal.Join("left join cicd_pipeline_commits cpr on cpr.repo_id = ? and cicd_tasks.pipeline_id = cpr.pipeline_id ", repoId),
			dal.Where("status=? ", devops.DONE))
	} else {

		cursor, err = db.Cursor(
			dal.From(&devops.CICDTask{}),
			dal.Join("left join cicd_pipeline_commits cpr on cpr.repo_id = ? and cicd_tasks.pipeline_id = cpr.pipeline_id ", repoId),
			dal.Where("status=? and SUBSTRING_INDEX(id, ':', 1) in ? ", devops.DONE, dataSource))
	}
	if err != nil {
		return err
	}

	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:    taskCtx,
			Params: DoraApiParams{
				// TODO
			},
			Table: "cicd_tasks",
		},
		InputRowType: reflect.TypeOf(devops.CICDTask{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			cicdTask := inputRow.(*devops.CICDTask)
			results := make([]interface{}, 0, 1)
			var EnvironmentVar string
			if productionNamePattern == "" {
				EnvironmentVar = devops.PRODUCTION
			} else {
				if productEnv := productionNameRegexp.FindString(cicdTask.Name); productEnv != "" {
					EnvironmentVar = devops.PRODUCTION
				}
			}

			// TODO: STAGE 2
			// if stagingEnv := stagingNameRegexp.FindString(cicdTask.Name); stagingEnv != "" {
			// 	EnvironmentVar = devops.STAGING
			// }
			// if testingEnv := testingNameRegexp.FindString(cicdTask.Name); testingEnv != "" {
			// 	EnvironmentVar = devops.TESTING
			// }

			cicdPipelineFilter := &devops.CICDTask{
				DomainEntity: cicdTask.DomainEntity,
				PipelineId:   cicdTask.PipelineId,
				Name:         cicdTask.Name,
				Type:         cicdTask.Type,
				Result:       cicdTask.Result,
				Status:       cicdTask.Status,
				DurationSec:  cicdTask.DurationSec,
				StartedDate:  cicdTask.StartedDate,
				FinishedDate: cicdTask.FinishedDate,
				Environment:  EnvironmentVar,
			}
			results = append(results, cicdPipelineFilter)
			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
