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
	"regexp"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
)

type JiraTransformationRule struct {
	common.Model               `mapstructure:"-"`
	ConnectionId               uint64          `mapstructure:"connectionId" json:"connectionId"`
	Name                       string          `mapstructure:"name" json:"name" gorm:"type:varchar(255);index:idx_name_jira,unique" validate:"required"`
	EpicKeyField               string          `mapstructure:"epicKeyField,omitempty" json:"epicKeyField" gorm:"type:varchar(255)"`
	StoryPointField            string          `mapstructure:"storyPointField,omitempty" json:"storyPointField" gorm:"type:varchar(255)"`
	RemotelinkCommitShaPattern string          `mapstructure:"remotelinkCommitShaPattern,omitempty" json:"remotelinkCommitShaPattern" gorm:"type:varchar(255)"`
	RemotelinkRepoPattern      json.RawMessage `mapstructure:"remotelinkRepoPattern,omitempty" json:"remotelinkRepoPattern"`
	TypeMappings               json.RawMessage `mapstructure:"typeMappings,omitempty" json:"typeMappings"`
}

func (r JiraTransformationRule) TableName() string {
	return "_tool_jira_transformation_rules"
}

func (r JiraTransformationRule) VerifyRegexp() errors.Error {
	var err error
	if r.RemotelinkCommitShaPattern != "" {
		_, err = regexp.Compile(r.RemotelinkCommitShaPattern)
		if err != nil {
			return errors.Convert(err)
		}
	}
	var repoPatterns []string
	err = json.Unmarshal(r.RemotelinkRepoPattern, &repoPatterns)
	if err != nil {
		return errors.Convert(err)
	}
	for _, pattern := range repoPatterns {
		_, err = regexp.Compile(pattern)
		if err != nil {
			return errors.Convert(err)
		}
	}
	return nil
}
