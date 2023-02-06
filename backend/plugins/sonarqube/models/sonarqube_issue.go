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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type SonarqubeIssue struct {
	common.NoPKModel
	ConnectionId uint64           `gorm:"primaryKey"`
	Key          string           `json:"key" gorm:"primaryKey"`
	BatchId      string           `json:"batchId" gorm:"type:varchar(100)"` // from collection time
	Rule         string           `json:"rule"`
	Severity     string           `json:"severity"`
	Component    string           `json:"component"`
	Project      string           `json:"project"`
	Line         int              `json:"line"`
	Status       string           `json:"status"`
	Message      string           `json:"message"`
	Debt         string           `json:"debt"`
	Effort       string           `json:"effort"`
	Author       string           `json:"author"`
	Hash         string           `json:"hash"`
	Tags         string           `json:"tags"`
	Type         string           `json:"type"`
	Scope        string           `json:"scope"`
	StartLine    int              `json:"startLine"`
	EndLine      int              `json:"endLine"`
	StartOffset  int              `json:"startOffset"`
	EndOffset    int              `json:"endOffset"`
	CreationDate *api.Iso8601Time `json:"creationDate"`
	UpdateDate   *api.Iso8601Time `json:"updateDate"`
}

func (SonarqubeIssue) TableName() string {
	return "_tool_sonarqube_issues"
}
