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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var _ core.SubTaskEntryPoint = ExtractDepartment

var ExtractDepartmentMeta = core.SubTaskMeta{
	Name:             "extractDepartment",
	EntryPoint:       ExtractDepartment,
	EnabledByDefault: true,
	Description:      "extract Zentao department",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ExtractDepartment(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ZentaoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProductId:    data.Options.ProductId,
				ExecutionId:  data.Options.ExecutionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_DEPARTMENT_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			department := &models.ZentaoDepartment{}
			err := json.Unmarshal(row.Data, department)
			if err != nil {
				return nil, errors.Default.WrapRaw(err)
			}
			department.ConnectionId = data.Options.ConnectionId
			results := make([]interface{}, 0)
			results = append(results, department)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
