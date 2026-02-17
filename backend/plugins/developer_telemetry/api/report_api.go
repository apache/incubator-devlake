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

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/developer_telemetry/models"
	"github.com/apache/incubator-devlake/server/api/shared"
)

type TelemetryReportRequest struct {
	Date      string           `json:"date" mapstructure:"date" validate:"required"`
	Developer string           `json:"developer" mapstructure:"developer" validate:"required"`
	Email     string           `json:"email" mapstructure:"email"`
	Name      string           `json:"name" mapstructure:"name"`
	Hostname  string           `json:"hostname" mapstructure:"hostname"`
	Metrics   TelemetryMetrics `json:"metrics" mapstructure:"metrics" validate:"required"`
}

type TelemetryMetrics struct {
	ActiveHours int            `json:"active_hours" mapstructure:"active_hours"`
	ToolsUsed   []string       `json:"tools_used" mapstructure:"tools_used"`
	Commands    map[string]int `json:"commands" mapstructure:"commands"`
	Projects    []string       `json:"projects" mapstructure:"projects"`
}

func PostReport(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	db := basicRes.GetDal()

	connectionIdStr := input.Params["connectionId"]
	connectionId, err := strconv.ParseUint(connectionIdStr, 10, 64)
	if err != nil {
		return nil, errors.BadInput.New("invalid connectionId")
	}

	connection := &models.DeveloperTelemetryConnection{}
	err = db.First(connection, dal.Where("id = ?", connectionId))
	if err != nil {
		if db.IsErrorNotFound(err) {
			return nil, errors.NotFound.New(fmt.Sprintf("connection %d not found", connectionId))
		}
		return nil, errors.Default.Wrap(err, "failed to find connection")
	}

	// TODO: Add authentication check
	// For now, authentication is disabled to test the flow
	// if connection.SecretToken != "" {
	//     // Need to find correct way to access request headers
	// }

	logger := basicRes.GetLogger()

	var report TelemetryReportRequest
	if err := api.Decode(input.Body, &report, nil); err != nil {
		return nil, errors.BadInput.Wrap(err, "failed to decode request body")
	}

	if report.Date == "" || report.Developer == "" {
		return nil, errors.BadInput.New("date and developer are required fields")
	}

	reportDate, parseErr := time.Parse("2006-01-02", report.Date)
	if parseErr != nil {
		return nil, errors.BadInput.Wrap(parseErr, "invalid date format, expected YYYY-MM-DD")
	}

	toolsUsedJSON, _ := json.Marshal(report.Metrics.ToolsUsed)
	projectsJSON, _ := json.Marshal(report.Metrics.Projects)
	commandsJSON, _ := json.Marshal(report.Metrics.Commands)

	metric := &models.DeveloperMetrics{
		ConnectionId:   connectionId,
		DeveloperId:    report.Developer,
		Date:           reportDate,
		Email:          report.Email,
		Name:           report.Name,
		Hostname:       report.Hostname,
		ActiveHours:    report.Metrics.ActiveHours,
		ToolsUsed:      string(toolsUsedJSON),
		ProjectContext: string(projectsJSON),
		CommandCounts:  string(commandsJSON),
		OsInfo:         "",
	}

	// Check if record exists using Count
	// Use formatted date string for comparison since the DB column is DATE type
	dateStr := reportDate.Format("2006-01-02")
	count, countErr := db.Count(
		dal.From(&models.DeveloperMetrics{}),
		dal.Where("connection_id = ? AND developer_id = ? AND date = ?",
			connectionId, report.Developer, dateStr))

	if countErr != nil {
		return nil, errors.Default.Wrap(countErr, "failed to check existing record")
	}

	var saveErr errors.Error
	if count == 0 {
		// Record doesn't exist, create new one
		saveErr = db.Create(metric)
	} else {
		// Record exists, update only non-primary-key columns
		saveErr = db.UpdateColumns(&models.DeveloperMetrics{}, []dal.DalSet{
			{ColumnName: "email", Value: metric.Email},
			{ColumnName: "name", Value: metric.Name},
			{ColumnName: "hostname", Value: metric.Hostname},
			{ColumnName: "active_hours", Value: metric.ActiveHours},
			{ColumnName: "tools_used", Value: metric.ToolsUsed},
			{ColumnName: "project_context", Value: metric.ProjectContext},
			{ColumnName: "command_counts", Value: metric.CommandCounts},
			{ColumnName: "os_info", Value: metric.OsInfo},
		}, dal.Where("connection_id = ? AND developer_id = ? AND date = ?",
			connectionId, report.Developer, dateStr))
	}

	if saveErr != nil {
		return nil, errors.Default.Wrap(saveErr, "failed to save developer metrics")
	}

	logger.Info("Successfully ingested telemetry data for developer=%s, date=%s, connection=%d",
		report.Developer, report.Date, connectionId)

	return &plugin.ApiResourceOutput{
		Body: shared.ApiBody{
			Success: true,
			Message: "telemetry data received successfully",
		},
		Status: http.StatusOK,
	}, nil
}
