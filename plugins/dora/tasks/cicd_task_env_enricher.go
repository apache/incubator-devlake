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
	EnabledByDefault: false,
	Description:      "calculate deployment frequency",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func EnrichTasksEnv(taskCtx core.SubTaskContext) (err errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*DoraTaskData)
	projectName := data.Options.ProjectName

	productionNamePattern := data.Options.ProductionPattern
	productionNameRegexp, err := errors.Convert01(regexp.Compile(productionNamePattern))
	if err != nil {
		return err
	}

	cursor, err := db.Cursor(
		dal.From(`cicd_tasks ct`),
		dal.Join("left join project_mapping pm on pm.row_id = ct.cicd_scope_id"),
		dal.Where(`pm.project_name = ? and pm.table = ?`, projectName, "cicd_scopes"),
	)

	if err != nil {
		return err
	}

	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: DoraApiParams{
				ProjectName: projectName,
			},
			Table: "cicd_tasks",
		},
		InputRowType: reflect.TypeOf(devops.CICDTask{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			cicdTask := inputRow.(*devops.CICDTask)
			if productionNamePattern == "" || productionNameRegexp.FindString(cicdTask.Name) != "" {
				cicdTask.Environment = devops.PRODUCTION
				return []interface{}{cicdTask}, nil
			}
			return nil, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
