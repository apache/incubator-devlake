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
	common.NoPKModel
	CommitFileID string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha    string `gorm:"type:varchar(40)"`
	FilePath     string `gorm:"type:varchar(255)"`
	Additions    int
	Deletions    int
}

func (CommitFile) TableName() string {
	return "commit_files"
}

type FileComponent struct {
	RepoId    string `gorm:"primaryKey;type:varchar(255)"`
	Component string `gorm:"primaryKey;type:varchar(255)"`
	PathRegex string `gorm:"type:varchar(255)"`
}

func (FileComponent) TableName() string {
	return "file_component"
}

type CommitfileComponent struct {
	common.NoPKModel
	CommitFileID string `gorm:"primaryKey;type:varchar(255)"`
	RepoId       string `gorm:"primaryKey;type:varchar(255)"`
	Component    string `gorm:"type:varchar(255)"`
}

func (CommitfileComponent) TableName() string {
	return "commitfile_component"
}
