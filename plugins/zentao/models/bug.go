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
	"time"
)

type ZentaoBugRes struct {
	ID             uint64     `json:"id"`
	Project        uint64     `json:"project"`
	Product        uint64     `json:"product"`
	Injection      int        `json:"injection"`
	Identify       int        `json:"identify"`
	Branch         int        `json:"branch"`
	Module         int        `json:"module"`
	Execution      uint64     `json:"execution"`
	Plan           int        `json:"plan"`
	Story          uint64     `json:"story"`
	StoryVersion   int        `json:"storyVersion"`
	Task           int        `json:"task"`
	ToTask         int        `json:"toTask"`
	ToStory        uint64     `json:"toStory"`
	Title          string     `json:"title"`
	Keywords       string     `json:"keywords"`
	Severity       int        `json:"severity"`
	Pri            int        `json:"pri"`
	Type           string     `json:"type"`
	Os             string     `json:"os"`
	Browser        string     `json:"browser"`
	Hardware       string     `json:"hardware"`
	Found          string     `json:"found"`
	Steps          string     `json:"steps"`
	Status         string     `json:"status"`
	SubStatus      string     `json:"subStatus"`
	Color          string     `json:"color"`
	Confirmed      int        `json:"confirmed"`
	ActivatedCount int        `json:"activatedCount"`
	ActivatedDate  *time.Time `json:"activatedDate"`
	FeedbackBy     string     `json:"feedbackBy"`
	NotifyEmail    string     `json:"notifyEmail"`
	OpenedBy       struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"openedBy"`
	OpenedDate  *time.Time `json:"openedDate"`
	OpenedBuild string     `json:"openedBuild"`
	AssignedTo  struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"assignedTo"`
	AssignedDate *time.Time `json:"assignedDate"`
	Deadline     string     `json:"deadline"`
	ResolvedBy   struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"resolvedBy"`
	Resolution    string     `json:"resolution"`
	ResolvedBuild string     `json:"resolvedBuild"`
	ResolvedDate  *time.Time `json:"resolvedDate"`
	ClosedBy      struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"closedBy"`
	ClosedDate   *time.Time `json:"closedDate"`
	DuplicateBug int        `json:"duplicateBug"`
	LinkBug      string     `json:"linkBug"`
	Feedback     int        `json:"feedback"`
	Result       int        `json:"result"`
	Repo         int        `json:"repo"`
	Mr           int        `json:"mr"`
	Entry        string     `json:"entry"`
	Lines        string     `json:"lines"`
	V1           string     `json:"v1"`
	V2           string     `json:"v2"`
	RepoType     string     `json:"repoType"`
	IssueKey     string     `json:"issueKey"`
	Testtask     int        `json:"testtask"`
	LastEditedBy struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"lastEditedBy"`
	LastEditedDate *time.Time `json:"lastEditedDate"`
	Deleted        bool       `json:"deleted"`
	PriOrder       string     `json:"priOrder"`
	SeverityOrder  int        `json:"severityOrder"`
	Needconfirm    bool       `json:"needconfirm"`
	StatusName     string     `json:"statusName"`
	ProductStatus  string     `json:"productStatus"`
}

type ZentaoBug struct {
	common.NoPKModel
	ConnectionId   uint64     `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ID             uint64     `json:"id" gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Project        uint64     `json:"project"`
	Product        uint64     `json:"product"`
	Injection      int        `json:"injection"`
	Identify       int        `json:"identify"`
	Branch         int        `json:"branch"`
	Module         int        `json:"module"`
	Execution      uint64     `json:"execution"`
	Plan           int        `json:"plan"`
	Story          uint64     `json:"story"`
	StoryVersion   int        `json:"storyVersion"`
	Task           int        `json:"task"`
	ToTask         int        `json:"toTask"`
	ToStory        uint64     `json:"toStory"`
	Title          string     `json:"title"`
	Keywords       string     `json:"keywords"`
	Severity       int        `json:"severity"`
	Pri            int        `json:"pri"`
	Type           string     `json:"type"`
	Os             string     `json:"os"`
	Browser        string     `json:"browser"`
	Hardware       string     `json:"hardware"`
	Found          string     `json:"found"`
	Steps          string     `json:"steps"`
	Status         string     `json:"status"`
	SubStatus      string     `json:"subStatus"`
	Color          string     `json:"color"`
	Confirmed      int        `json:"confirmed"`
	ActivatedCount int        `json:"activatedCount"`
	ActivatedDate  *time.Time `json:"activatedDate"`
	FeedbackBy     string     `json:"feedbackBy"`
	NotifyEmail    string     `json:"notifyEmail"`
	OpenedById     uint64
	OpenedByName   string
	OpenedDate     *time.Time `json:"openedDate"`
	OpenedBuild    string     `json:"openedBuild"`
	AssignedToId   uint64
	AssignedToName string
	AssignedDate   *time.Time `json:"assignedDate"`
	Deadline       string     `json:"deadline"`
	ResolvedById   uint64
	Resolution     string     `json:"resolution"`
	ResolvedBuild  string     `json:"resolvedBuild"`
	ResolvedDate   *time.Time `json:"resolvedDate"`
	ClosedById     uint64
	ClosedDate     *time.Time `json:"closedDate"`
	DuplicateBug   int        `json:"duplicateBug"`
	LinkBug        string     `json:"linkBug"`
	Feedback       int        `json:"feedback"`
	Result         int        `json:"result"`
	Repo           int        `json:"repo"`
	Mr             int        `json:"mr"`
	Entry          string     `json:"entry"`
	Lines          string     `json:"lines"`
	V1             string     `json:"v1"`
	V2             string     `json:"v2"`
	RepoType       string     `json:"repoType"`
	IssueKey       string     `json:"issueKey"`
	Testtask       int        `json:"testtask"`
	LastEditedById uint64
	LastEditedDate *time.Time `json:"lastEditedDate"`
	Deleted        bool       `json:"deleted"`
	PriOrder       string     `json:"priOrder"`
	SeverityOrder  int        `json:"severityOrder"`
	Needconfirm    bool       `json:"needconfirm"`
	StatusName     string     `json:"statusName"`
	ProductStatus  string     `json:"productStatus"`
}

func (ZentaoBug) TableName() string {
	return "_tool_zentao_bugs"
}
