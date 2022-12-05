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

type JiraBoard struct {
	common.NoPKModel     `json:"-" mapstructure:"-"`
	ConnectionId         uint64 `json:"connectionId" mapstructure:"connectionId" gorm:"primaryKey"`
	BoardId              uint64 `json:"boardId" mapstructure:"boardId" gorm:"primaryKey"`
	TransformationRuleId uint64 `json:"transformationRuleId,omitempty" mapstructure:"transformationRuleId"`
	ProjectId            uint   `json:"projectId" mapstructure:"projectId"`
	Name                 string `json:"name" mapstructure:"name" gorm:"type:varchar(255)"`
	Self                 string `json:"self" mapstructure:"self" gorm:"type:varchar(255)"`
	Type                 string `json:"type" mapstructure:"type" gorm:"type:varchar(100)"`
}

func (JiraBoard) TableName() string {
	return "_tool_jira_boards"
}
