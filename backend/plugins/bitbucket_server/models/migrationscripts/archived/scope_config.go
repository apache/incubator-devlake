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

package archived

import (
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"gorm.io/datatypes"
)

type BitbucketServerScopeConfig struct {
	archived.ScopeConfig `mapstructure:",squash" json:",inline" gorm:"embedded"`
	ConnectionId         uint64 `json:"connectionId" gorm:"index" validate:"required" mapstructure:"connectionId,omitempty"`
	Name                 string `mapstructure:"name" json:"name" gorm:"type:varchar(255);uniqueIndex" validate:"required"`
	PrType               string `mapstructure:"prType,omitempty" json:"prType" gorm:"type:varchar(255)"`
	PrComponent          string `mapstructure:"prComponent,omitempty" json:"prComponent" gorm:"type:varchar(255)"`
	PrBodyClosePattern   string `mapstructure:"prBodyClosePattern,omitempty" json:"prBodyClosePattern" gorm:"type:varchar(255)"`

	// DeploymentPattern  string            `mapstructure:"deploymentPattern,omitempty" json:"deploymentPattern" gorm:"type:varchar(255)"`
	// ProductionPattern  string            `mapstructure:"productionPattern,omitempty" json:"productionPattern" gorm:"type:varchar(255)"`
	Refdiff datatypes.JSONMap `mapstructure:"refdiff,omitempty" json:"refdiff" swaggertype:"object" format:"json"`

	// a string array, split by `,`.
}

func (BitbucketServerScopeConfig) TableName() string {
	return "_tool_bitbucket_server_scope_configs"
}
