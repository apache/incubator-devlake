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

var _ core.SubTaskEntryPoint = ExtractProducts

var ExtractProductMeta = core.SubTaskMeta{
	Name:             "extractProducts",
	EntryPoint:       ExtractProducts,
	EnabledByDefault: true,
	Description:      "extract Zentao products",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func ExtractProducts(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ZentaoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ExecutionId:  data.Options.ExecutionId,
				ProductId:    data.Options.ProductId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_PRODUCT_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			res := &models.ZentaoProductRes{}
			err := json.Unmarshal(row.Data, res)
			if err != nil {
				return nil, errors.Default.Wrap(err, "error reading endpoint response by Zentao product extractor")
			}
			product := &models.ZentaoProduct{
				ConnectionId:   data.Options.ConnectionId,
				Id:             uint64(res.ID),
				Program:        res.Program,
				Name:           res.Name,
				Code:           res.Code,
				Bind:           res.Bind,
				Line:           res.Line,
				Type:           res.Type,
				Status:         res.Status,
				SubStatus:      res.SubStatus,
				Description:    res.Description,
				POId:           res.PO.ID,
				QDId:           res.QD.ID,
				RDId:           res.RD.ID,
				Acl:            res.Acl,
				Reviewer:       res.Reviewer,
				CreatedById:    res.CreatedBy.ID,
				CreatedDate:    res.CreatedDate,
				CreatedVersion: res.CreatedVersion,
				OrderIn:        res.OrderIn,
				Deleted:        res.Deleted,
				Plans:          res.Plans,
				Releases:       res.Releases,
				Builds:         res.Builds,
				Cases:          res.Cases,
				Projects:       res.Projects,
				Executions:     res.Executions,
				Bugs:           res.Bugs,
				Docs:           res.Docs,
				Progress:       res.Progress,
				CaseReview:     res.CaseReview,
			}
			results := make([]interface{}, 0)
			results = append(results, product)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
