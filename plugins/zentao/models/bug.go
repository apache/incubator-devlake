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
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type ZentaoBug struct {
	common.NoPKModel
	ConnectionId   uint64   `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ID             int      `json:"id" gorm:"primaryKey"`
	Project        int      `json:"project"`
	Product        int      `json:"product"`
	Injection      int      `json:"injection"`
	Identify       int      `json:"identify"`
	Branch         int      `json:"branch"`
	Module         int      `json:"module"`
	Execution      int      `json:"execution"`
	Plan           int      `json:"plan"`
	Story          int      `json:"story"`
	StoryVersion   int      `json:"storyVersion"`
	Task           int      `json:"task"`
	ToTask         int      `json:"toTask"`
	ToStory        int      `json:"toStory"`
	Title          string   `json:"title"`
	Keywords       string   `json:"keywords"`
	Severity       int      `json:"severity"`
	Pri            int      `json:"pri"`
	Type           string   `json:"type"`
	Os             string   `json:"os"`
	Browser        string   `json:"browser"`
	Hardware       string   `json:"hardware"`
	Found          string   `json:"found"`
	Steps          string   `json:"steps"`
	Status         string   `json:"status"`
	SubStatus      string   `json:"subStatus"`
	Color          string   `json:"color"`
	Confirmed      int      `json:"confirmed"`
	ActivatedCount int      `json:"activatedCount"`
	ActivatedDate  string   `json:"activatedDate"`
	FeedbackBy     string   `json:"feedbackBy"`
	NotifyEmail    string   `json:"notifyEmail"`
	Mailto         []Mailto `json:"mailto" gorm:"-:all"`
	OpenedBy       `json:"openedBy"`
	OpenedDate     *helper.Iso8601Time `json:"openedDate"`
	OpenedBuild    string              `json:"openedBuild"`
	AssignedTo     `json:"assignedTo"`
	AssignedDate   *helper.Iso8601Time `json:"assignedDate"`
	Deadline       string              `json:"deadline"`
	ResolvedBy     `json:"resolvedBy"`
	Resolution     string              `json:"resolution"`
	ResolvedBuild  string              `json:"resolvedBuild"`
	ResolvedDate   *helper.Iso8601Time `json:"resolvedDate"`
	ClosedBy       `json:"closedBy"`
	ClosedDate     *helper.Iso8601Time `json:"closedDate"`
	DuplicateBug   int                 `json:"duplicateBug"`
	LinkBug        string              `json:"linkBug"`
	Case           int                 `json:"case"`
	CaseVersion    int                 `json:"caseVersion"`
	Feedback       int                 `json:"feedback"`
	Result         int                 `json:"result"`
	Repo           int                 `json:"repo"`
	Mr             int                 `json:"mr"`
	Entry          string              `json:"entry"`
	Lines          string              `json:"lines"`
	V1             string              `json:"v1"`
	V2             string              `json:"v2"`
	RepoType       string              `json:"repoType"`
	IssueKey       string              `json:"issueKey"`
	Testtask       int                 `json:"testtask"`
	LastEditedBy   `json:"lastEditedBy"`
	LastEditedDate *helper.Iso8601Time `json:"lastEditedDate"`
	Deleted        bool                `json:"deleted"`
	PriOrder       string              `json:"priOrder"`
	SeverityOrder  int                 `json:"severityOrder"`
	Needconfirm    bool                `json:"needconfirm"`
	StatusName     string              `json:"statusName"`
	ProductStatus  string              `json:"productStatus"`
}
type ResolvedBy struct {
	ResolvedByID       int    `json:"id"`
	ResolvedByAccount  string `json:"account"`
	ResolvedByAvatar   string `json:"avatar"`
	ResolvedByRealname string `json:"realname"`
}
type Mailto struct {
	MailtoID       int    `json:"id"`
	MailtoAccount  string `json:"account"`
	MailtoAvatar   string `json:"avatar"`
	MailtoRealname string `json:"realname"`
}

func (ZentaoBug) TableName() string {
	return "_tool_zentao_bugs"
}
