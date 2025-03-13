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

// Deployment represents the entire JSON structure
type Deployment struct {
	common.NoPKModel `swaggerignore:"true" json:"-" mapstructure:"-"`

	ConnectionId    uint64 `gorm:"primaryKey"`
	ProjectId       string `gorm:"type:varchar(255)"`
	Name            string `gorm:"type:varchar(255)"`
	GeneratedName   string `gorm:"type:varchar(255)"`
	Namespace       string `gorm:"type:varchar(255)"`
	UID             string `gorm:"type:varchar(255)"`
	ResourceVersion string `gorm:"type:varchar(255)"`
	Result          string `gorm:"type:varchar(255)"`
	CreationDate    string `gorm:"type:varchar(255)"`
	StartedAt       string `gorm:"type:varchar(255)"`
	FinishedAt      string `gorm:"type:varchar(255)"`
	CommitSha       string `gorm:"type:varchar(255)"`
	RefName         string `gorm:"type:varchar(255)"`
	RepoUrl         string `gorm:"type:varchar(255)"`
	DurationSec     int64  `gorm:"type:integer(12)"`
}

func (Deployment) TableName() string {
	return "_tool_argo_api_deployments"
}
