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
	"time"
)

type Commit struct {
	NoPKModel
	Sha            string `json:"sha" gorm:"primaryKey;type:varchar(40);comment:commit hash"`
	Additions      int    `gorm:"comment:Added lines of code"`
	Deletions      int    `gorm:"comment:Deleted lines of code"`
	DevEq          int    `gorm:"comment:Merico developer equivalent from analysis engine"`
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

type CommitFile struct {
	NoPKModel
	CommitSha string `gorm:"primaryKey;type:varchar(40)"`
	FilePath  string `gorm:"primaryKey;type:varchar(255)"`
	Additions int
	Deletions int
}

type CommitParent struct {
	CommitSha       string `json:"commitSha" gorm:"primaryKey;type:varchar(40);comment:commit hash"`
	ParentCommitSha string `json:"parentCommitSha" gorm:"primaryKey;type:varchar(40);comment:parent commit hash"`
}
