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

	"github.com/apache/incubator-devlake/core/models/common"
)

const (
	BLUEPRINT_MODE_NORMAL   = "NORMAL"
	BLUEPRINT_MODE_ADVANCED = "ADVANCED"
)

// @Description CronConfig
type Blueprint struct {
	Name         string                 `json:"name" validate:"required"`
	ProjectName  string                 `json:"projectName" gorm:"type:varchar(255)"`
	Mode         string                 `json:"mode" gorm:"varchar(20)" validate:"required,oneof=NORMAL ADVANCED"`
	Plan         PipelinePlan           `json:"plan" gorm:"serializer:encdec"`
	Enable       bool                   `json:"enable"`
	CronConfig   string                 `json:"cronConfig" format:"* * * * *" example:"0 0 * * 1"`
	IsManual     bool                   `json:"isManual"`
	SkipOnFail   bool                   `json:"skipOnFail"`
	BeforePlan   PipelinePlan           `json:"beforePlan" gorm:"serializer:encdec"`
	AfterPlan    PipelinePlan           `json:"afterPlan" gorm:"serializer:encdec"`
	TimeAfter    *time.Time             `json:"timeAfter"`
	Labels       []string               `json:"labels" gorm:"-"`
	Connections  []*BlueprintConnection `json:"connections" gorm:"-"`
	common.Model `swaggerignore:"true"`
}

func (Blueprint) TableName() string {
	return "_devlake_blueprints"
}

type BlueprintLabel struct {
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	BlueprintId uint64    `json:"blueprint_id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"primaryKey;index"`
}

func (BlueprintLabel) TableName() string {
	return "_devlake_blueprint_labels"
}

type BlueprintConnection struct {
	BlueprintId  uint64            `json:"-" gorm:"primaryKey" validate:"required"`
	PluginName   string            `json:"pluginName" gorm:"primaryKey;type:varchar(255)" validate:"required"`
	ConnectionId uint64            `json:"connectionId" gorm:"primaryKey" validate:"required"`
	Scopes       []*BlueprintScope `json:"scopes" gorm:"-"`
}

func (BlueprintConnection) TableName() string {
	return "_devlake_blueprint_connections"
}

type BlueprintScope struct {
	BlueprintId  uint64 `json:"-" gorm:"primaryKey" validate:"required"`
	PluginName   string `json:"-" gorm:"primaryKey;type:varchar(255)" validate:"required"`
	ConnectionId uint64 `json:"-" gorm:"primaryKey" validate:"required"`
	ScopeId      string `json:"scopeId" gorm:"primaryKey;type:varchar(255)" validate:"required"`
}

func (BlueprintScope) TableName() string {
	return "_devlake_blueprint_scopes"
}

type BlueprintSyncPolicy struct {
	SkipOnFail bool       `json:"skipOnFail"`
	TimeAfter  *time.Time `json:"timeAfter"`
}
