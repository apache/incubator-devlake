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
	"time"

	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type BitbucketServerRepo struct {
	archived.NoPKModel
	ConnectionId  uint64     `json:"connectionId" gorm:"primaryKey" validate:"required" mapstructure:"connectionId,omitempty"`
	ScopeConfigId uint64     `json:"scopeConfigId,omitempty" mapstructure:"scopeConfigId,omitempty"`
	BitbucketId   string     `json:"bitbucketId" gorm:"primaryKey;type:varchar(255)" validate:"required" mapstructure:"bitbucketId"`
	Name          string     `json:"name" gorm:"type:varchar(255)" mapstructure:"name,omitempty"`
	HTMLUrl       string     `json:"HTMLUrl" gorm:"type:varchar(255)" mapstructure:"HTMLUrl,omitempty"`
	Description   string     `json:"description" mapstructure:"description,omitempty"`
	CloneUrl      string     `json:"cloneUrl" gorm:"type:varchar(255)" mapstructure:"cloneUrl,omitempty"`
	CreatedDate   *time.Time `json:"createdDate" mapstructure:"-"`
	UpdatedDate   *time.Time `json:"updatedDate" mapstructure:"-"`
}

func (BitbucketServerRepo) TableName() string {
	return "_tool_bitbucket_server_repos"
}
