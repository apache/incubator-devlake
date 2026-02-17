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

package migrationscripts

import (
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*addInitTables)(nil)

type addInitTables struct{}

type developerTelemetryConnection20240211 struct {
	ID               uint64    `gorm:"primaryKey;type:BIGINT  NOT NULL AUTO_INCREMENT" json:"id"`
	Name             string    `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	Endpoint         string    `json:"endpoint"`
	SecretToken      string    `mapstructure:"secretToken" json:"secretToken" gorm:"serializer:encdec"`
	Proxy            string    `json:"proxy"`
	RateLimitPerHour int       `comment:"api request rate limit per hour" json:"rateLimitPerHour"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

func (developerTelemetryConnection20240211) TableName() string {
	return "_tool_developer_telemetry_connections"
}

type developerMetrics20240211 struct {
	ConnectionId   uint64    `gorm:"primaryKey;type:BIGINT" json:"connection_id"`
	DeveloperId    string    `gorm:"primaryKey;type:varchar(255)" json:"developer_id"`
	Date           string    `gorm:"primaryKey;type:date" json:"date"`
	Email          string    `gorm:"type:varchar(255);index" json:"email"`
	Name           string    `gorm:"type:varchar(255)" json:"name"`
	Hostname       string    `gorm:"type:varchar(255)" json:"hostname"`
	ActiveHours    int       `json:"active_hours"`
	ToolsUsed      string    `gorm:"type:text" json:"tools_used"`
	ProjectContext string    `gorm:"type:text" json:"project_context"`
	CommandCounts  string    `gorm:"type:text" json:"command_counts"`
	OsInfo         string    `gorm:"type:varchar(255)" json:"os_info"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (developerMetrics20240211) TableName() string {
	return "_tool_developer_metrics"
}

func (*addInitTables) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.AutoMigrate(&developerTelemetryConnection20240211{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&developerMetrics20240211{})
	if err != nil {
		return err
	}
	return nil
}

func (*addInitTables) Version() uint64 {
	return 20240211000001
}

func (*addInitTables) Name() string {
	return "developer_telemetry init schemas"
}
