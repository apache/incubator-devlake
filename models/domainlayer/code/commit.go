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

package code

import (
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer"

	"github.com/apache/incubator-devlake/models/common"
)

type Commit struct {
	common.NoPKModel
	Sha            string `json:"sha" gorm:"primaryKey;type:varchar(40);comment:commit hash"`
	Additions      int    `json:"additions" gorm:"comment:Added lines of code"`
	Deletions      int    `json:"deletions" gorm:"comment:Deleted lines of code"`
	DevEq          int    `json:"deveq" gorm:"comment:Merico developer equivalent from analysis engine"`
	Message        string
	AuthorName     string `gorm:"type:varchar(255)"`
	AuthorEmail    string `gorm:"type:varchar(255)"`
	AuthoredDate   time.Time
	AuthorId       string `gorm:"type:varchar(255)"`
	CommitterName  string `gorm:"type:varchar(255)"`
	CommitterEmail string `gorm:"type:varchar(255)"`
	CommittedDate  time.Time
	CommitterId    string `gorm:"index;type:varchar(255)"`
}

func (Commit) TableName() string {
	return "commits"
}

type CommitFile struct {
	domainlayer.DomainEntity
	CommitSha string `gorm:"type:varchar(40)"`
	FilePath  string `gorm:"type:text"`
	Additions int
	Deletions int
}

func (CommitFile) TableName() string {
	return "commit_files"
}

type CommitFileComponent struct {
	common.NoPKModel
	CommitFileId  string `gorm:"primaryKey;type:varchar(255)"`
	ComponentName string `gorm:"type:varchar(255)"`
}

func (CommitFileComponent) TableName() string {
	return "commit_file_components"
}
