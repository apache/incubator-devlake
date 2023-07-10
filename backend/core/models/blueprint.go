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
	"encoding/json"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
)

const (
	BLUEPRINT_MODE_NORMAL   = "NORMAL"
	BLUEPRINT_MODE_ADVANCED = "ADVANCED"
)

type (
	Blueprint struct {
		Name        string          `json:"name" validate:"required"`
		ProjectName string          `json:"projectName" gorm:"type:varchar(255)"`
		Mode        string          `json:"mode" gorm:"varchar(20)" validate:"required,oneof=NORMAL ADVANCED"`
		Plan        json.RawMessage `json:"plan" gorm:"serializer:encdec"`
		Enable      bool            `json:"enable"`
		//please check this https://crontab.guru/ for detail
		CronConfig   string             `json:"cronConfig" format:"* * * * *" example:"0 0 * * 1"`
		IsManual     bool               `json:"isManual"`
		SkipOnFail   bool               `json:"skipOnFail"`
		Labels       []string           `json:"labels" gorm:"-"`
		Settings     *BlueprintSettings `json:"settings" gorm:"-"`
		common.Model `swaggerignore:"true"`
	}
	BlueprintSettings struct {
		common.Model `swaggerignore:"true"`
		Version      string
		BlueprintId  uint64
		TimeAfter    *time.Time             `json:"timeAfter"`
		Connections  []*BlueprintConnection `json:"connections" gorm:"-" validate:"required"`
		BeforePlan   json.RawMessage        `json:"before_plan"`
		AfterPlan    json.RawMessage        `json:"after_plan"`
	}
	BlueprintConnection struct {
		common.Model `swaggerignore:"true"`
		BlueprintId  uint64
		ConnectionId uint64
		SettingsId   uint64
		Plugin       string
		Scopes       []*BlueprintScope `json:"scopes" gorm:"-" validate:"required"`
	}
	BlueprintScope struct {
		common.Model `swaggerignore:"true"`
		Name         string
		ScopeId      string
		ConnectionId uint64
		BlueprintId  uint64
	}
)

// UnmarshalPlan unmarshals Plan in JSON to strong-typed plugin.PipelinePlan
func (bp *Blueprint) UnmarshalPlan() (plugin.PipelinePlan, errors.Error) {
	var plan plugin.PipelinePlan
	err := errors.Convert(json.Unmarshal(bp.Plan, &plan))
	if err != nil {
		return nil, errors.Default.Wrap(err, `unmarshal plan fail`)
	}
	return plan, nil
}

// GetScopes Gets all the scopes for a given connection for this blueprint. Returns an empty slice if none found.
func (bp *Blueprint) GetScopes(connectionId uint64, pluginName string) ([]*BlueprintScope, errors.Error) {
	conns := bp.Settings.Connections
	visited := map[string]any{}
	var result []*BlueprintScope
	for _, conn := range conns {
		if conn.ConnectionId != connectionId || conn.Plugin != pluginName {
			continue
		}
		for _, scope := range conn.Scopes {
			if _, ok := visited[scope.ScopeId]; !ok {
				result = append(result, scope)
				visited[scope.ScopeId] = true
			}
		}
	}
	return result, nil
}

func (*Blueprint) TableName() string {
	return "_devlake_blueprints"
}

func (*BlueprintSettings) TableName() string {
	return "_devlake_blueprint_settings"
}

func (*BlueprintConnection) TableName() string {
	return "_devlake_blueprint_connections"
}

func (*BlueprintScope) TableName() string {
	return "_devlake_blueprint_scopes"
}

type DbBlueprintLabel struct {
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	BlueprintId uint64    `json:"blueprint_id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"primaryKey;index"`
}

func (*DbBlueprintLabel) TableName() string {
	return "_devlake_blueprint_labels"
}
