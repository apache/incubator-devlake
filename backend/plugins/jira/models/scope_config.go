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
	"regexp"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
)

type StatusMapping struct {
	StandardStatus string `json:"standardStatus"`
}

type StatusMappings map[string]StatusMapping

type TypeMapping struct {
	StandardType   string         `json:"standardType"`
	StatusMappings StatusMappings `json:"statusMappings"`
}

type CommitUrlPattern struct {
	Pattern string `json:"pattern"`
	Regex   string `json:"regex"`
}

type JiraScopeConfig struct {
	common.ScopeConfig         `mapstructure:",squash" json:",inline" gorm:"embedded"`
	ConnectionId               uint64                 `mapstructure:"connectionId" json:"connectionId"`
	Name                       string                 `mapstructure:"name" json:"name" gorm:"type:varchar(255);index:idx_name_jira,unique" validate:"required"`
	EpicKeyField               string                 `mapstructure:"epicKeyField,omitempty" json:"epicKeyField" gorm:"type:varchar(255)"`
	StoryPointField            string                 `mapstructure:"storyPointField,omitempty" json:"storyPointField" gorm:"type:varchar(255)"`
	RemotelinkCommitShaPattern string                 `mapstructure:"remotelinkCommitShaPattern,omitempty" json:"remotelinkCommitShaPattern" gorm:"type:varchar(255)"`
	RemotelinkRepoPattern      []CommitUrlPattern     `mapstructure:"remotelinkRepoPattern,omitempty" json:"remotelinkRepoPattern" gorm:"type:json;serializer:json"`
	TypeMappings               map[string]TypeMapping `mapstructure:"typeMappings,omitempty" json:"typeMappings" gorm:"type:json;serializer:json"`
	ApplicationType            string                 `mapstructure:"applicationType,omitempty" json:"applicationType" gorm:"type:varchar(255)"`
}

func (r *JiraScopeConfig) SetConnectionId(c *JiraScopeConfig, connectionId uint64) {
	c.ConnectionId = connectionId
	c.ScopeConfig.ConnectionId = connectionId
}

func (r *JiraScopeConfig) Validate() errors.Error {
	var err error
	if r.RemotelinkCommitShaPattern != "" {
		_, err = regexp.Compile(r.RemotelinkCommitShaPattern)
		if err != nil {
			return errors.Convert(err)
		}
	}
	for _, pattern := range r.RemotelinkRepoPattern {
		if pattern.Regex == "" {
			return errors.BadInput.New("empty regex in remotelinkRepoPattern")
		}
		_, err = regexp.Compile(pattern.Regex)
		if err != nil {
			return errors.Convert(err)
		}
	}
	return nil
}

func (r JiraScopeConfig) TableName() string {
	return "_tool_jira_scope_configs"
}
