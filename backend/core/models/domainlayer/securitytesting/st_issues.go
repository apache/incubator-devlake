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

package securitytesting

import (
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type StIssue struct {
	domainlayer.DomainEntity
	Rule                     string           `json:"rule" gorm:"type:varchar(255)"`
	Severity                 string           `json:"severity" gorm:"type:varchar(255)"`
	Component                string           `json:"component" gorm:"type:varchar(255)"`
	ProjectKey               string           `gorm:"index;type:varchar(255)"`          //domain project key
	BatchId                  string           `json:"batchId" gorm:"type:varchar(100)"` // from collection time
	Line                     int              `json:"line"`
	Status                   string           `json:"status" gorm:"type:varchar(255)"`
	Message                  string           `json:"message"`
	Debt                     string           `json:"debt" gorm:"type:varchar(255)"`
	Effort                   string           `json:"effort" gorm:"type:varchar(255)"`
	CommitAuthorEmail        string           `json:"author" gorm:"type:varchar(255)"`
	Assignee                 string           `json:"assignee" gorm:"type:varchar(255)"`
	Hash                     string           `json:"hash" gorm:"type:varchar(255)"`
	Tags                     string           `json:"tags" gorm:"type:varchar(255)"`
	Type                     string           `json:"type" gorm:"type:varchar(255)"`
	Scope                    string           `json:"scope" gorm:"type:varchar(255)"`
	StartLine                int              `json:"startLine"`
	EndLine                  int              `json:"endLine"`
	StartOffset              int              `json:"startOffset"`
	EndOffset                int              `json:"endOffset"`
	VulnerabilityProbability string           `gorm:"type:varchar(100)"`
	SecurityCategory         string           `gorm:"type:varchar(100)"`
	CreationDate             *api.Iso8601Time `json:"creationDate"`
	UpdateDate               *api.Iso8601Time `json:"updateDate"`
}

func (StIssue) TableName() string {
	return "st_issues"
}
