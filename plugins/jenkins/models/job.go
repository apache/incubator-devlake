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
	"github.com/apache/incubator-devlake/models/common"
)

// JenkinsJob db entity for jenkins job
type JenkinsJob struct {
	ConnectionId         uint64 `gorm:"primaryKey" mapstructure:"connectionId,omitempty" json:"connectionId"`
	FullName             string `gorm:"primaryKey;type:varchar(255)" mapstructure:"jobFullName" json:"jobFullName"` // "path1/path2/job name"
	TransformationRuleId uint64 `mapstructure:"transformationRules,omitempty" json:"transformationRuleId,omitempty"`
	Name                 string `gorm:"index;type:varchar(255)" mapstructure:"-,omitempty" json:"-"` // "job name"
	Path                 string `gorm:"index;type:varchar(511)" mapstructure:"-,omitempty" json:"-"` // "job/path1/job/path2"
	Class                string `gorm:"type:varchar(255)" mapstructure:"class,omitempty" json:"class"`
	Color                string `gorm:"type:varchar(255)" mapstructure:"color,omitempty" json:"color"`
	Base                 string `gorm:"type:varchar(255)" mapstructure:"base,omitempty" json:"base"`
	Url                  string `mapstructure:"url,omitempty" json:"url"`
	Description          string `mapstructure:"description,omitempty" json:"description"`
	PrimaryView          string `gorm:"type:varchar(255)" mapstructure:"primaryView,omitempty" json:"primaryView"`
	common.NoPKModel     `json:"-" mapstructure:"-"`
}

func (JenkinsJob) TableName() string {
	return "_tool_jenkins_jobs"
}
