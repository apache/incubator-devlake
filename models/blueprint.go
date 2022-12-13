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

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/core"
)

const (
	BLUEPRINT_MODE_NORMAL   = "NORMAL"
	BLUEPRINT_MODE_ADVANCED = "ADVANCED"
)

// @Description CronConfig
type Blueprint struct {
	Name        string          `json:"name" validate:"required"`
	ProjectName string          `json:"projectName" gorm:"type:varchar(255)"`
	Mode        string          `json:"mode" gorm:"varchar(20)" validate:"required,oneof=NORMAL ADVANCED"`
	Plan        json.RawMessage `json:"plan"`
	Enable      bool            `json:"enable"`
	//please check this https://crontab.guru/ for detail
	CronConfig   string          `json:"cronConfig" format:"* * * * *" example:"0 0 * * 1"`
	IsManual     bool            `json:"isManual"`
	SkipOnFail   bool            `json:"skipOnFail"`
	Labels       []string        `json:"labels"`
	Settings     json.RawMessage `json:"settings" swaggertype:"array,string" example:"please check api: /blueprints/<PLUGIN_NAME>/blueprint-setting"`
	common.Model `swaggerignore:"true"`
}

type BlueprintSettings struct {
	Version          string          `json:"version" validate:"required,semver,oneof=1.0.0"`
	CreatedDateAfter *time.Time      `json:"createdDateAfter"`
	Connections      json.RawMessage `json:"connections" validate:"required"`
	BeforePlan       json.RawMessage `json:"before_plan"`
	AfterPlan        json.RawMessage `json:"after_plan"`
}

// UnmarshalPlan unmarshals Plan in JSON to strong-typed core.PipelinePlan
func (bp *Blueprint) UnmarshalPlan() (core.PipelinePlan, errors.Error) {
	var plan core.PipelinePlan
	err := errors.Convert(json.Unmarshal(bp.Plan, &plan))
	if err != nil {
		return nil, errors.Convert(err)
	}
	return plan, nil
}

// @Description CronConfig
type DbBlueprint struct {
	Name        string `json:"name" validate:"required"`
	ProjectName string `json:"projectName" gorm:"type:varchar(255)"`
	Mode        string `json:"mode" gorm:"varchar(20)" validate:"required,oneof=NORMAL ADVANCED"`
	Plan        string `json:"plan" encrypt:"yes"`
	Enable      bool   `json:"enable"`
	//please check this https://crontab.guru/ for detail
	CronConfig       string     `json:"cronConfig" format:"* * * * *" example:"0 0 * * 1"`
	IsManual         bool       `json:"isManual"`
	SkipOnFail       bool       `json:"skipOnFail"`
	CreatedDateAfter *time.Time `json:"createdDateAfter"`
	Settings         string     `json:"settings" encrypt:"yes" swaggertype:"array,string" example:"please check api: /blueprints/<PLUGIN_NAME>/blueprint-setting"`
	common.Model     `swaggerignore:"true"`

	Labels []DbBlueprintLabel `json:"-" gorm:"-"`
}

func (DbBlueprint) TableName() string {
	return "_devlake_blueprints"
}

type DbBlueprintLabel struct {
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	BlueprintId uint64    `json:"blueprint_id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"primaryKey;index"`
}

func (DbBlueprintLabel) TableName() string {
	return "_devlake_blueprint_labels"
}
