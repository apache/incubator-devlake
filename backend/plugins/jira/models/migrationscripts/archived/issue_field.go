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

import "github.com/apache/incubator-devlake/core/models/migrationscripts/archived"

type JiraIssueField struct {
	archived.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey"`
	BoardId      uint64 `gorm:"primaryKey"`

	ID               string `json:"id" gorm:"primaryKey"`
	Name             string `json:"name"`
	Custom           bool   `json:"custom"`
	Orderable        bool   `json:"orderable"`
	Navigable        bool   `json:"navigable"`
	Searchable       bool   `json:"searchable"`
	SchemaType       string `json:"schema_type"`
	SchemaItems      string `json:"schema_items"`
	SchemaCustom     string `json:"schema_custom"`
	SchemaCustomID   int    `json:"schema_custom_id"`
	ScheCustomSystem string `json:"sche_custom_system"`
}

func (JiraIssueField) TableName() string {
	return "_tool_jira_issue_fields"
}
