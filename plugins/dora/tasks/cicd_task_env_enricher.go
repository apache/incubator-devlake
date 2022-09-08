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
	"fmt"
	"reflect"
	"regexp"
	"runtime/debug"

	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/dora/api"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var EnrichTaskEnvMeta = core.SubTaskMeta{
	Name:             "EnrichTaskEnv",
	EntryPoint:       EnrichTasksEnv,
	EnabledByDefault: true,
	Description:      "calculate deployment frequency",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func EnrichTasksEnv(taskCtx core.SubTaskContext) (err error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*DoraTaskData)

	var environmentRegex *regexp.Regexp
	environmentPattern := data.Options.TransformationRules
	if len(environmentPattern) > 0 {
		environmentRegex, err = regexp.Compile(environmentPattern)
		if err != nil {
			return fmt.Errorf("regexp Compile environmentRegex failed:[%s] stack:[%s]", err.Error(), debug.Stack())
		}
	} else {
		fmt.Println("NO SETTING")
		environmentPattern = "deploy" // default
		environmentRegex, err = regexp.Compile("deploy")
		if err != nil {
			return fmt.Errorf("regexp Compile environmentRegex failed:[%s] stack:[%s]", err.Error(), debug.Stack())
		}
	}
	fmt.Println(environmentRegex)

	cursor, err := db.Cursor(
		dal.From(&devops.CICDTask{}),
		dal.Where("status=?", devops.DONE))
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
			Table: api.RAW_DEPLOYMENTS_TABLE,
		},
		InputRowType: reflect.TypeOf(devops.CICDTask{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			cicdTask := inputRow.(*devops.CICDTask)
			results := make([]interface{}, 0, 1)
			if deployTask := environmentRegex.FindString(cicdTask.Name); deployTask == "" {
				return nil, nil
			}
			cicdPipelineFilter := &devops.CICDTask{
				DomainEntity: cicdTask.DomainEntity,
				Name:         cicdTask.Name,
				Result:       cicdTask.Result,
				Status:       cicdTask.Status,
				Type:         "DEPLOY",
				DurationSec:  cicdTask.DurationSec,
				StartedDate:  cicdTask.StartedDate,
				FinishedDate: cicdTask.FinishedDate,
				Environment:  environmentPattern,
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
