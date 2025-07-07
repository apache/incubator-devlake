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
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	testmoModels "github.com/apache/incubator-devlake/plugins/testmo/models"
)

var ConvertProjectsMeta = plugin.SubTaskMeta{
	Name:             "convertProjects",
	EntryPoint:       ConvertProjects,
	EnabledByDefault: true,
	Description:      "Convert tool layer table testmo_projects into domain layer table projects",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}

func ConvertProjects(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*TestmoTaskData)
	db := taskCtx.GetDal()

	cursor, err := db.Cursor(dal.From(&testmoModels.TestmoProject{}), dal.Where("connection_id = ? AND id = ?", data.Options.ConnectionId, data.Options.ProjectId))
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TestmoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_PROJECT_TABLE,
		},
		InputRowType: reflect.TypeOf(testmoModels.TestmoProject{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			project := inputRow.(*testmoModels.TestmoProject)

			domainProject := &models.Project{
				BaseProject: models.BaseProject{
					Name:        project.Name,
					Description: "",
				},
			}

			return []interface{}{domainProject}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
