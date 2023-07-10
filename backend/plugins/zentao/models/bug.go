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
	"bytes"
	"encoding/json"
	"github.com/apache/incubator-devlake/core/models/common"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type ApiAccount struct {
	ID       int64  `json:"id"`
	Account  string `json:"account"`
	Avatar   string `json:"avatar"`
	Realname string `json:"realname"`
}

func (a *ApiAccount) UnmarshalJSON(data []byte) error {
	var dst struct {
		ID       int64  `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	}
	data = bytes.TrimSpace(data)
	if string(data) == "null" {
		a = nil
	}
	if len(data) > 1 && data[0] == '"' && data[len(data)-1] == '"' {
		dst.Account = string(data[1 : len(data)-1])
	} else {
		err := json.Unmarshal(data, &dst)
		if err != nil {
			return err
		}
	}
	*a = dst
	return nil
}

type ZentaoBugRes struct {
	ID             int64               `json:"id"`
	Project        int64               `json:"project"`
	Product        int64               `json:"product"`
	Injection      int                 `json:"injection"`
	Identify       int                 `json:"identify"`
	Branch         int                 `json:"branch"`
	Module         int                 `json:"module"`
	Execution      int64               `json:"execution"`
	Plan           int                 `json:"plan"`
	Story          int64               `json:"story"`
	StoryVersion   int                 `json:"storyVersion"`
	Task           int                 `json:"task"`
	ToTask         int                 `json:"toTask"`
	ToStory        int64               `json:"toStory"`
	Title          string              `json:"title"`
	Keywords       string              `json:"keywords"`
	Severity       int                 `json:"severity"`
	Pri            int                 `json:"pri"`
	Type           string              `json:"type"`
	Os             string              `json:"os"`
	Browser        string              `json:"browser"`
	Hardware       string              `json:"hardware"`
	Found          string              `json:"found"`
	Steps          string              `json:"steps"`
	Status         string              `json:"status"`
	SubStatus      string              `json:"subStatus"`
	Color          string              `json:"color"`
	Confirmed      int                 `json:"confirmed"`
	ActivatedCount int                 `json:"activatedCount"`
	ActivatedDate  *helper.Iso8601Time `json:"activatedDate"`
	FeedbackBy     string              `json:"feedbackBy"`
	NotifyEmail    string              `json:"notifyEmail"`
	OpenedBy       *ApiAccount         `json:"openedBy"`
	OpenedDate     *helper.Iso8601Time `json:"openedDate"`
	OpenedBuild    string              `json:"openedBuild"`
	AssignedTo     *ApiAccount         `json:"assignedTo"`
	AssignedDate   *helper.Iso8601Time `json:"assignedDate"`
	Deadline       string              `json:"deadline"`
	ResolvedBy     *ApiAccount         `json:"resolvedBy"`
	Resolution     string              `json:"resolution"`
	ResolvedBuild  string              `json:"resolvedBuild"`
	ResolvedDate   *helper.Iso8601Time `json:"resolvedDate"`
	ClosedBy       *ApiAccount         `json:"closedBy"`
	ClosedDate     *helper.Iso8601Time `json:"closedDate"`
	DuplicateBug   int                 `json:"duplicateBug"`
	LinkBug        string              `json:"linkBug"`
	Feedback       int                 `json:"feedback"`
	Result         int                 `json:"result"`
	Repo           int                 `json:"repo"`
	Mr             int                 `json:"mr"`
	Entry          string              `json:"entry"`
	NumOfLine      string              `json:"lines"`
	V1             string              `json:"v1"`
	V2             string              `json:"v2"`
	RepoType       string              `json:"repoType"`
	IssueKey       string              `json:"issueKey"`
	Testtask       int                 `json:"testtask"`
	LastEditedBy   *ApiAccount         `json:"lastEditedBy"`
	LastEditedDate *helper.Iso8601Time `json:"lastEditedDate"`
	Deleted        bool                `json:"deleted"`
	PriOrder       string              `json:"priOrder"`
	SeverityOrder  int                 `json:"severityOrder"`
	Needconfirm    bool                `json:"needconfirm"`
	StatusName     string              `json:"statusName"`
	ProductStatus  string              `json:"productStatus"`
}

type ZentaoBug struct {
	common.NoPKModel
	ConnectionId   uint64              `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ID             int64               `json:"id" gorm:"primaryKey;type:BIGINT  NOT NULL;autoIncrement:false"`
	Project        int64               `json:"project"`
	Product        int64               `json:"product"`
	Injection      int                 `json:"injection"`
	Identify       int                 `json:"identify"`
	Branch         int                 `json:"branch"`
	Module         int                 `json:"module"`
	Execution      int64               `json:"execution"`
	Plan           int                 `json:"plan"`
	Story          int64               `json:"story"`
	StoryVersion   int                 `json:"storyVersion"`
	Task           int                 `json:"task"`
	ToTask         int                 `json:"toTask"`
	ToStory        int64               `json:"toStory"`
	Title          string              `json:"title"`
	Keywords       string              `json:"keywords"`
	Severity       int                 `json:"severity"`
	Pri            int                 `json:"pri"`
	Type           string              `json:"type"`
	Os             string              `json:"os"`
	Browser        string              `json:"browser"`
	Hardware       string              `json:"hardware"`
	Found          string              `json:"found"`
	Steps          string              `json:"steps"`
	Status         string              `json:"status"`
	SubStatus      string              `json:"subStatus"`
	Color          string              `json:"color"`
	Confirmed      int                 `json:"confirmed"`
	ActivatedCount int                 `json:"activatedCount"`
	ActivatedDate  *helper.Iso8601Time `json:"activatedDate"`
	FeedbackBy     string              `json:"feedbackBy"`
	NotifyEmail    string              `json:"notifyEmail"`
	OpenedById     int64
	OpenedByName   string
	OpenedDate     *helper.Iso8601Time `json:"openedDate"`
	OpenedBuild    string              `json:"openedBuild"`
	AssignedToId   int64
	AssignedToName string
	AssignedDate   *helper.Iso8601Time `json:"assignedDate"`
	Deadline       string              `json:"deadline"`
	ResolvedById   int64
	Resolution     string              `json:"resolution"`
	ResolvedBuild  string              `json:"resolvedBuild"`
	ResolvedDate   *helper.Iso8601Time `json:"resolvedDate"`
	ClosedById     int64
	ClosedDate     *helper.Iso8601Time `json:"closedDate"`
	DuplicateBug   int                 `json:"duplicateBug"`
	LinkBug        string              `json:"linkBug"`
	Feedback       int                 `json:"feedback"`
	Result         int                 `json:"result"`
	Repo           int                 `json:"repo"`
	Mr             int                 `json:"mr"`
	Entry          string              `json:"entry"`
	NumOfLine      string              `json:"lines"`
	V1             string              `json:"v1"`
	V2             string              `json:"v2"`
	RepoType       string              `json:"repoType"`
	IssueKey       string              `json:"issueKey"`
	Testtask       int                 `json:"testtask"`
	LastEditedById int64
	LastEditedDate *helper.Iso8601Time `json:"lastEditedDate"`
	Deleted        bool                `json:"deleted"`
	PriOrder       string              `json:"priOrder"`
	SeverityOrder  int                 `json:"severityOrder"`
	Needconfirm    bool                `json:"needconfirm"`
	StatusName     string              `json:"statusName"`
	ProductStatus  string              `json:"productStatus"`
	Url            string              `json:"url"`
	StdStatus      string              `json:"stdStatus" gorm:"type:varchar(20)"`
	StdType        string              `json:"stdType" gorm:"type:varchar(20)"`
}

func (ZentaoBug) TableName() string {
	return "_tool_zentao_bugs"
}
