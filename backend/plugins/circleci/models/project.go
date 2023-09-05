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

import "github.com/apache/incubator-devlake/core/models/common"

type CircleciProject struct {
	ConnectionId     uint64          `gorm:"primaryKey;type:BIGINT"`
	Id               string          `gorm:"primaryKey;type:varchar(100)" json:"id"`
	ProjectSlug      string          `gorm:"uniqueIndex;type:varchar(255)" json:"project_slug"`
	Slug             string          `gorm:"type:varchar(255)" json:"slug"`
	Name             string          `gorm:"type:varchar(255)" json:"name"`
	OrganizationName string          `gorm:"type:varchar(255)" json:"organization_name"`
	OrganizationSlug string          `gorm:"type:varchar(255)" json:"organization_slug"`
	OrganizationId   string          `gorm:"type:varchar(100)" json:"organization_id"`
	VcsInfo          CircleciVcsInfo `gorm:"serializer:json;type:text" json:"vcs_info"`

	common.NoPKModel
}

type CircleciVcsInfo struct {
	VcsUrl        string `json:"vcs_url"`
	Provider      string `json:"provider"`
	DefaultBranch string `json:"default_branch"`
}

func (CircleciProject) TableName() string {
	return "_tool_circleci_projects"
}
