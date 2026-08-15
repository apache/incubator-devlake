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

// GithubIssueFieldValue holds one issue field value for one issue.
//
// Issue fields are organization-level structured metadata that GitHub made generally
// available on 2026-07-02. Unlike labels they are typed, mutually exclusive within a
// field, and shared across every repository in the organization.
//
// Value carries a queryable text form of the value regardless of DataType:
//   - single_select: the selected option name
//   - multi_select:  the selected option names, comma separated, in API order
//   - text:          the text
//   - number:        the number, formatted without a trailing ".0" when integral
//   - date:          the date as returned by the API (YYYY-MM-DD)
//
// RawValue keeps the original JSON so a value that does not fit the text form
// (a multi_select array, for instance) is still recoverable downstream.
type GithubIssueFieldValue struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	IssueId      int    `gorm:"primaryKey;comment:GitHub issue id, matches _tool_github_issues.github_id"`
	FieldId      int    `gorm:"primaryKey"`
	FieldName    string `gorm:"type:varchar(255);index"`
	DataType     string `gorm:"type:varchar(50)"`
	Value        string `gorm:"type:text"`
	RawValue     string `gorm:"type:text"`
	OptionColor  string `gorm:"type:varchar(50)"`
	common.NoPKModel
}

func (GithubIssueFieldValue) TableName() string {
	return "_tool_github_issue_field_values"
}
