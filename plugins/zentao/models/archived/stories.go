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
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type ZentaoStories struct {
	archived.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ExecutionId  uint64 `json:"execution_id"`
	Project      int    `json:"project"`
	Product      int    `json:"product"`
	Branch       int    `json:"branch"`
	Story        int    `json:"story"`
	Version      int    `json:"version"`
	Order        int    `json:"order"`
	ID           int    `json:"id" gorm:"primaryKey;type:BIGINT  NOT NULL" `
	Vision       string `json:"vision"`
	Parent       int    `json:"parent"`
	Module       int    `json:"module"`
	Plan         string `json:"plan"`
	Source       string `json:"source"`
	SourceNote   string `json:"sourceNote"`
	FromBug      int    `json:"fromBug"`
	Feedback     int    `json:"feedback"`
	Title        string `json:"title"`
	Keywords     string `json:"keywords"`
	Type         string `json:"type"`
	Category     string `json:"category"`
	Pri          int    `json:"pri"`
	Estimate     int    `json:"estimate"`
	Status       string `json:"status"`
	SubStatus    string `json:"subStatus"`
	Color        string `json:"color"`
	Stage        string `json:"stage"`
	StagedBy     string `json:"stagedBy"`
	//Mailto           []interface{} `json:"mailto"`
	Lib            int `json:"lib"`
	FromStory      int `json:"fromStory"`
	FromVersion    int `json:"fromVersion"`
	OpenedBy       `json:"openedBy"`
	OpenedDate     *helper.Iso8601Time `json:"openedDate"`
	AssignedTo     `json:"assignedTo"`
	AssignedDate   *helper.Iso8601Time `json:"assignedDate"`
	ApprovedDate   string              `json:"approvedDate"`
	LastEditedBy   `json:"lastEditedBy"`
	LastEditedDate *helper.Iso8601Time `json:"lastEditedDate"`
	ChangedBy      string              `json:"changedBy"`
	ChangedDate    string              `json:"changedDate"`
	//ReviewedBy       interface{} `json:"reviewedBy"`
	ReviewedDate     *helper.Iso8601Time `json:"reviewedDate"`
	ClosedBy         `json:"closedBy"`
	ClosedDate       *helper.Iso8601Time `json:"closedDate"`
	ClosedReason     string              `json:"closedReason"`
	ActivatedDate    string              `json:"activatedDate"`
	ToBug            int                 `json:"toBug"`
	ChildStories     string              `json:"childStories"`
	LinkStories      string              `json:"linkStories"`
	LinkRequirements string              `json:"linkRequirements"`
	DuplicateStory   int                 `json:"duplicateStory"`
	StoryChanged     string              `json:"storyChanged"`
	FeedbackBy       string              `json:"feedbackBy"`
	NotifyEmail      string              `json:"notifyEmail"`
	URChanged        string              `json:"URChanged"`
	Deleted          bool                `json:"deleted"`
	PriOrder         string              `json:"priOrder"`
	ProductType      string              `json:"productType"`
	PlanTitle        string              `json:"planTitle"`
	ProductStatus    string              `json:"productStatus"`
}

type AssignedTo struct {
	ID       int    `json:"id"`
	Account  string `json:"account"`
	Avatar   string `json:"avatar"`
	Realname string `json:"realname"`
}

func (ZentaoStories) TableName() string {
	return "_tool_zentao_stories"
}
