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

package models

import (
	"time"
)

// DeveloperMetrics represents the tool layer table for developer telemetry data
type DeveloperMetrics struct {
	ConnectionId   uint64    `gorm:"primaryKey;type:BIGINT;column:connection_id" json:"connection_id"`
	DeveloperId    string    `gorm:"primaryKey;type:varchar(255);column:developer_id" json:"developer_id"`
	Date           time.Time `gorm:"primaryKey;type:date;column:date" json:"date"`
	Email          string    `gorm:"type:varchar(255);index;column:email" json:"email"`
	Name           string    `gorm:"type:varchar(255);column:name" json:"name"`
	Hostname       string    `gorm:"type:varchar(255);column:hostname" json:"hostname"`
	ActiveHours    int       `gorm:"column:active_hours" json:"active_hours"`
	ToolsUsed      string    `gorm:"type:text;column:tools_used" json:"tools_used"`           // JSON array stored as text
	ProjectContext string    `gorm:"type:text;column:project_context" json:"project_context"` // JSON array stored as text
	CommandCounts  string    `gorm:"type:text;column:command_counts" json:"command_counts"`   // JSON object stored as text
	OsInfo         string    `gorm:"type:varchar(255);column:os_info" json:"os_info"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (DeveloperMetrics) TableName() string {
	return "_tool_developer_metrics"
}
