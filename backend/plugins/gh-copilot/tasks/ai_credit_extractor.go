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
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gh-copilot/models"
)

// aiCreditUsageRecord represents a single usage item from the AI credit usage API.
type aiCreditUsageRecord struct {
	Product          string  `json:"product"`
	Sku              string  `json:"sku"`
	Model            string  `json:"model"`
	UnitType         string  `json:"unitType"`
	PricePerUnit     float64 `json:"pricePerUnit"`
	GrossQuantity    float64 `json:"grossQuantity"`
	DiscountQuantity float64 `json:"discountQuantity"`
	NetQuantity      float64 `json:"netQuantity"`
	GrossAmount      float64 `json:"grossAmount"`
	DiscountAmount   float64 `json:"discountAmount"`
	NetAmount        float64 `json:"netAmount"`
}

// aiCreditResponseWrapper represents the wrapper around the API response containing time period and usage items.
type aiCreditResponseWrapper struct {
	TimePeriod struct {
		Year  int `json:"year"`
		Month int `json:"month"`
		Day   int `json:"day"`
	} `json:"timePeriod"`
	Enterprise   string `json:"enterprise"`
	Organization string `json:"organization"`
	User         string `json:"user"`
	Product      string `json:"product"`
	Model        string `json:"model"`
	CostCenter   struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"costCenter"`
}

// ExtractAiCreditUsage parses AI credit usage records into the appropriate model tables.
func ExtractAiCreditUsage(taskCtx plugin.SubTaskContext) errors.Error {
	data, ok := taskCtx.TaskContext().GetData().(*GhCopilotTaskData)
	if !ok {
		return errors.Default.New("task data is not GhCopilotTaskData")
	}
	connection := data.Connection
	connection.Normalize()

	extractor, err := helper.NewApiExtractor(taskCtx, rawAiCreditUsageTable)
	if err != nil {
		return err
	}

	err = extractor.Extract(
		rawAiCreditUsageTable,
		func(row *helper.RawData) ([]interface{}, errors.Error) {
			// Parse raw data
			var record aiCreditUsageRecord
			err := json.Unmarshal(row.Data, &record)
			if err != nil {
				return nil, errors.Convert(err)
			}

			// Extract wrapper info from row context
			var wrapper aiCreditResponseWrapper
			if scopeAttr, ok := row.Params["scope"]; ok {
				if connection.HasEnterprise() {
					wrapper.Enterprise = scopeAttr.(string)
				} else if connection.Organization != "" {
					wrapper.Organization = scopeAttr.(string)
				}
			}
			wrapper.Product = record.Product
			wrapper.Model = record.Model

			// Parse date from row (need to extract from response context)
			// For now, use current date - this should be enhanced to parse from API response
			now := time.Now().UTC()
			wrapper.TimePeriod.Year = now.Year()
			wrapper.TimePeriod.Month = int(now.Month())
			wrapper.TimePeriod.Day = now.Day()

			var results []interface{}

			// Route to appropriate table based on connection type
			if connection.HasEnterprise() {
				toolRecord := &models.GhCopilotEnterpriseAiCreditUsage{
					ConnectionId:     data.Connection.ID,
					ScopeId:          data.Options.ScopeId,
					Year:             wrapper.TimePeriod.Year,
					Month:            wrapper.TimePeriod.Month,
					Day:              wrapper.TimePeriod.Day,
					Enterprise:       wrapper.Enterprise,
					Model:            record.Model,
					Organization:     wrapper.Organization,
					User:             wrapper.User,
					Product:          record.Product,
					CostCenterId:     wrapper.CostCenter.Id,
					CostCenterName:   wrapper.CostCenter.Name,
					GrossQuantity:    record.GrossQuantity,
					DiscountQuantity: record.DiscountQuantity,
					NetQuantity:      record.NetQuantity,
					PricePerUnit:     record.PricePerUnit,
					GrossAmount:      record.GrossAmount,
					DiscountAmount:   record.DiscountAmount,
					NetAmount:        record.NetAmount,
				}
				results = append(results, toolRecord)
			} else if connection.Organization != "" {
				toolRecord := &models.GhCopilotOrgAiCreditUsage{
					ConnectionId:     data.Connection.ID,
					ScopeId:          data.Options.ScopeId,
					Year:             wrapper.TimePeriod.Year,
					Month:            wrapper.TimePeriod.Month,
					Day:              wrapper.TimePeriod.Day,
					Organization:     wrapper.Organization,
					Model:            record.Model,
					User:             wrapper.User,
					Product:          record.Product,
					GrossQuantity:    record.GrossQuantity,
					DiscountQuantity: record.DiscountQuantity,
					NetQuantity:      record.NetQuantity,
					PricePerUnit:     record.PricePerUnit,
					GrossAmount:      record.GrossAmount,
					DiscountAmount:   record.DiscountAmount,
					NetAmount:        record.NetAmount,
				}
				results = append(results, toolRecord)
			} else {
				// User-level credits
				toolRecord := &models.GhCopilotUserAiCreditUsage{
					ConnectionId:     data.Connection.ID,
					ScopeId:          data.Options.ScopeId,
					Year:             wrapper.TimePeriod.Year,
					Month:            wrapper.TimePeriod.Month,
					Day:              wrapper.TimePeriod.Day,
					User:             connection.GetEmail(), // Use authenticated user
					Model:            record.Model,
					Product:          record.Product,
					GrossQuantity:    record.GrossQuantity,
					DiscountQuantity: record.DiscountQuantity,
					NetQuantity:      record.NetQuantity,
					PricePerUnit:     record.PricePerUnit,
					GrossAmount:      record.GrossAmount,
					DiscountAmount:   record.DiscountAmount,
					NetAmount:        record.NetAmount,
				}
				results = append(results, toolRecord)
			}

			return results, nil
		},
	)

	return err
}
