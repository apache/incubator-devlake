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
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	bambooModels "github.com/apache/incubator-devlake/plugins/bamboo/models"
)

const RAW_PLAN_TABLE = "bamboo_plan"

var ConvertPlansMeta = plugin.SubTaskMeta{
	Name:             "convertplans",
	EntryPoint:       ConvertPlans,
	EnabledByDefault: true,
	Description:      "Convert tool layer table bamboo_plans into  domain layer table plans",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func ConvertPlans(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PLAN_TABLE)
	cursor, err := db.Cursor(dal.From(bambooModels.BambooPlan{}),
		dal.Where("connection_id = ? and plan_key = ?", data.Options.ConnectionId, data.Options.PlanKey))
	if err != nil {
		return err
	}
	defer cursor.Close()

	planIdGen := didgen.NewDomainIdGenerator(&bambooModels.BambooPlan{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(bambooModels.BambooPlan{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			bambooPlan := inputRow.(*bambooModels.BambooPlan)
			domainPlan := &devops.CicdScope{
				DomainEntity: domainlayer.DomainEntity{Id: planIdGen.Generate(data.Options.ConnectionId, bambooPlan.PlanKey)},
				Name:         bambooPlan.Name,
				Description:  bambooPlan.Description,
				Url:          bambooPlan.Href,
			}
			return []interface{}{
				domainPlan,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
