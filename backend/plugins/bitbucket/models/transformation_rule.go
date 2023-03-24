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
	"github.com/apache/incubator-devlake/core/models/common"
	"gorm.io/datatypes"
)

type BitbucketTransformationRule struct {
	common.Model      `mapstructure:"-"`
	ConnectionId      uint64            `mapstructure:"connectionId" json:"connectionId"`
	Name              string            `mapstructure:"name" json:"name" gorm:"type:varchar(255);index:idx_name_github,unique" validate:"required"`
	DeploymentPattern string            `mapstructure:"deploymentPattern,omitempty" json:"deploymentPattern" gorm:"type:varchar(255)"`
	ProductionPattern string            `mapstructure:"productionPattern,omitempty" json:"productionPattern" gorm:"type:varchar(255)"`
	Refdiff           datatypes.JSONMap `mapstructure:"refdiff,omitempty" json:"refdiff" swaggertype:"object" format:"json"`

	// a string array, split by `,`.
	IssueStatusTodo       string `mapstructure:"issueStatusTodo,omitempty" json:"issueStatusTodo" gorm:"type:varchar(255)"`
	IssueStatusInProgress string `mapstructure:"issueStatusInProgress,omitempty" json:"issueStatusInProgress" gorm:"type:varchar(255)"`
	IssueStatusDone       string `mapstructure:"issueStatusDone,omitempty" json:"issueStatusDone" gorm:"type:varchar(255)"`
	IssueStatusOther      string `mapstructure:"issueStatusOther,omitempty" json:"issueStatusOther" gorm:"type:varchar(255)"`
}

func (BitbucketTransformationRule) TableName() string {
	return "_tool_bitbucket_transformation_rules"
}
