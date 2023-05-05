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
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

const RAW_PRODUCT_TABLE = "zentao_api_products"

var _ plugin.SubTaskEntryPoint = ConvertProducts

var ConvertProductMeta = plugin.SubTaskMeta{
	Name:             "convertProducts",
	EntryPoint:       ConvertProducts,
	EnabledByDefault: true,
	Description:      "convert Zentao products",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertProducts(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	db := taskCtx.GetDal()
	boardIdGen := didgen.NewDomainIdGenerator(&models.ZentaoProduct{})
	cursor, err := db.Cursor(
		dal.From(&models.ZentaoProduct{}),
		dal.Where(`id = ? and connection_id = ?`, data.Options.ProductId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()
	convertor, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.ZentaoProduct{}),
		Input:        cursor,
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ZentaoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProductId:    data.Options.ProductId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_PRODUCT_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			toolProduct := inputRow.(*models.ZentaoProduct)

			domainBoard := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{
					Id: boardIdGen.Generate(toolProduct.ConnectionId, toolProduct.Id),
				},
				Name:        toolProduct.Name,
				Description: toolProduct.Description,
				CreatedDate: toolProduct.CreatedDate.ToNullableTime(),
				Type:        toolProduct.Type + "/" + toolProduct.ProductType,
			}
			results := make([]interface{}, 0)
			results = append(results, domainBoard)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return convertor.Execute()
}
