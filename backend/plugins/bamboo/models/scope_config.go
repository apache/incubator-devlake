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
)

const ENV_NAME_PATTERN = "ENV_NAME_PATTERN"

type BambooScopeConfig struct {
	common.ScopeConfig `mapstructure:",squash" json:",inline" gorm:"embedded"`
	ConnectionId       uint64           `mapstructure:"connectionId" json:"connectionId"`
	Name               string           `gorm:"type:varchar(255);index:idx_name_gitlab,unique" validate:"required" mapstructure:"name" json:"name"`
	RepoMap            map[string][]int `json:"repoMap" gorm:"type:json;serializer:json"` // should be {realRepoName: [bamboo_repoId, ...]}
	DeploymentPattern  string           `mapstructure:"deploymentPattern,omitempty" json:"deploymentPattern" gorm:"type:varchar(255)"`
	ProductionPattern  string           `mapstructure:"productionPattern,omitempty" json:"productionPattern" gorm:"type:varchar(255)"`
	EnvNamePattern     string           `mapstructure:"envNamePattern,omitempty" json:"envNamePattern" gorm:"type:varchar(255)"`
}

func (BambooScopeConfig) TableName() string {
	return "_tool_bamboo_scope_configs"
}

func (cfg *BambooScopeConfig) SetConnectionId(c *BambooScopeConfig, connectionId uint64) {
	c.ConnectionId = connectionId
	c.ScopeConfig.ConnectionId = connectionId
}

func (cfg *BambooScopeConfig) GetRepoIdMap() map[int]string {
	repoMap := make(map[int]string)
	for k, v := range cfg.RepoMap {
		for _, id := range v {
			repoMap[id] = k
		}
	}
	return repoMap
}
