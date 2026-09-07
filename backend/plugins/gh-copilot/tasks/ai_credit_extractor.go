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
	"time"

	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
	helper "github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/plugins/gh-copilot/models"
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

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:   taskCtx,
			Table: rawAiCreditUsageTable,
			Options: copilotRawParams{
				ConnectionId: data.Options.ConnectionId,
				ScopeId:      data.Options.ScopeId,
				Organization: connection.Organization,
				Endpoint:     connection.Endpoint,
			},
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			// Parse raw data
			var record aiCreditUsageRecord
			err := json.Unmarshal(row.Data, &record)
			if err != nil {
				return nil, errors.Convert(err)
			}

			// Extract wrapper info from row context
			var wrapper aiCreditResponseWrapper
			if connection.HasEnterprise() {
				wrapper.Enterprise = connection.Enterprise
			} else if connection.Organization != "" {
				wrapper.Organization = connection.Organization
			}
			wrapper.Product = record.Product
			wrapper.Model = record.Model

			// Derive the time period from the collector's day input so records are
			// deterministic and aligned with the requested billing day, rather than
			// depending on the extraction-time clock.
			var input dayInput
			if len(row.Input) > 0 {
				if unmErr := json.Unmarshal(row.Input, &input); unmErr != nil {
					return nil, errors.Convert(unmErr)
				}
			}
			day, parseErr := time.Parse("2006-01-02", input.Day)
			if parseErr != nil {
				return nil, errors.Convert(parseErr)
			}
			wrapper.TimePeriod.Year = day.Year()
			wrapper.TimePeriod.Month = int(day.Month())
			wrapper.TimePeriod.Day = day.Day()

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
					User:             connection.Name, // Use connection display name
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
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}
