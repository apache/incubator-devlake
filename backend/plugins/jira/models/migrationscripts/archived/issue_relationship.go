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
)

type JiraIssueRelationship struct {
	archived.NoPKModel
	ConnectionId    uint64 `gorm:"primaryKey"`
	IssueId         uint64 `gorm:"primarykey"`
	IssueKey        string `gorm:"type:varchar(255)"` // e.g. DEV-1
	TypeId          uint64 // e.g. 10001
	TypeName        string `gorm:"type:varchar(255)"` // e.g. Blocks
	Inward          string `gorm:"type:varchar(255)"` // e.g. blocks
	Outward         string `gorm:"type:varchar(255)"` // e.g. is blocked by
	InwardIssueId   uint64 // e.g. 116566
	InwardIssueKey  string `gorm:"type:varchar(255)"` // e.g. DEV-2
	OutwardIssueId  uint64 // e.g. 116567
	OutwardIssueKey string `gorm:"type:varchar(255)"` // e.g. DEV-3
}

func (JiraIssueRelationship) TableName() string {
	return "_tool_jira_issue_relationships"
}
