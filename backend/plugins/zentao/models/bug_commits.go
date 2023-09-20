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

type ZentaoBugCommitsRes struct {
	ID         int    `json:"id" gorm:"primaryKey;type:BIGINT  NOT NULL;autoIncrement:false"`
	ObjectType string `json:"objectType"`
	ObjectID   int    `json:"objectID"`
	Product    string `json:"product"`
	Project    int    `json:"project"`
	Execution  int    `json:"execution"`
	Actor      string `json:"actor"`
	Action     string `json:"action"`
	Date       string `json:"date"`
	Comment    string `json:"comment"`
	Extra      string `json:"extra"`
	Read       string `json:"read"`
	Vision     string `json:"vision"`
	Efforted   int    `json:"efforted"`
	Desc       string `json:"desc"`
}

type ZentaoBugCommit struct {
	common.NoPKModel
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

type ZentaoBugRepoCommitsRes struct {
	Title string `json:"title"`
	Log   struct {
		Revision  string `json:"revision"`
		Committer string `json:"committer"`
		Time      string `json:"time"`
		Comment   string `json:"comment"`
		Commit    string `json:"commit"`
	} `json:"log"`
	Repo struct {
		ID                 string `json:"id"`
		Product            string `json:"product"`
		Projects           string `json:"projects"`
		Name               string `json:"name"`
		Path               string `json:"path"`
		Prefix             string `json:"prefix"`
		Encoding           string `json:"encoding"`
		Scm                string `json:"SCM"`
		Client             string `json:"client"`
		ServiceHost        string `json:"serviceHost"`
		ServiceProject     string `json:"serviceProject"`
		Commits            string `json:"commits"`
		Account            string `json:"account"`
		Password           string `json:"password"`
		Encrypt            string `json:"encrypt"`
		ACL                any    `json:"acl"`
		Synced             string `json:"synced"`
		LastSync           string `json:"lastSync"`
		Desc               string `json:"desc"`
		Extra              string `json:"extra"`
		PreMerge           string `json:"preMerge"`
		Job                string `json:"job"`
		FileServerURL      any    `json:"fileServerUrl"`
		FileServerAccount  string `json:"fileServerAccount"`
		FileServerPassword string `json:"fileServerPassword"`
		Deleted            string `json:"deleted"`
		CodePath           string `json:"codePath"`
		GitService         string `json:"gitService"`
		Project            string `json:"project"`
	} `json:"repo"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	RepoID      string `json:"repoID"`
	BranchID    bool   `json:"branchID"`
	ObjectID    string `json:"objectID"`
	Revision    string `json:"revision"`
	ParentDir   string `json:"parentDir"`
	OldRevision string `json:"oldRevision"`
	PreAndNext  struct {
		Pre  string `json:"pre"`
		Next string `json:"next"`
	} `json:"preAndNext"`
	Pager any `json:"pager"`
}

type ZentaoBugRepoCommit struct {
	common.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Product      int64  `json:"product"`
	Project      int64  `json:"project"`
	IssueId      string `gorm:"primaryKey;type:varchar(255)"` // the bug id
	RepoUrl      string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha    string `gorm:"primaryKey;type:varchar(255)"`
}

func (ZentaoBugRepoCommit) TableName() string {
	return "_tool_zentao_bug_repo_commits"
}
