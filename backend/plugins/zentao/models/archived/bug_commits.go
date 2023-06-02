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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type ZentaoBugCommit struct {
	archived.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ID           int    `json:"id" gorm:"primaryKey;type:BIGINT  NOT NULL;autoIncrement:false"`
	ObjectType   string `json:"objectType"`
	ObjectID     int    `json:"objectID"`
	Product      int64  `json:"product"`
	Project      int64  `json:"project"`
	Execution    int    `json:"execution"`
	Actor        string `json:"actor"`
	Action       string `json:"action"`
	Date         string `json:"date"`
	Comment      string `json:"comment"`
	Extra        string `json:"extra"`
	Host         string `json:"host"`         //the host part of extra
	RepoRevision string `json:"repoRevision"` // the repoRevisionJson part of extra
	ActionRead   string `json:"actionRead"`
	Vision       string `json:"vision"`
	Efforted     int    `json:"efforted"`
	ActionDesc   string `json:"cctionDesc"`
}

func (ZentaoBugCommit) TableName() string {
	return "_tool_zentao_bug_commits"
}

type ZentaoBugRepoCommit  struct {
	common.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Product      int64  `json:"product"`
	Project      int64  `json:"project"`
	IssueId      string `gorm:"primaryKey;type:varchar(255)"` // the bug id
	RepoUrl      string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha    string `gorm:"primaryKey;type:varchar(255)"`
}

func (ZentaoBugRepoCommit ) TableName() string {
	return "_tool_zentao_bug_repo_commits"
}
