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

type ZentaoStoryRes struct {
	ID               uint64              `json:"id"`
	Vision           string              `json:"vision"`
	Parent           uint64              `json:"parent"`
	Product          uint64              `json:"product"`
	Branch           int                 `json:"branch"`
	Module           int                 `json:"module"`
	Plan             string              `json:"plan"`
	Source           string              `json:"source"`
	SourceNote       string              `json:"sourceNote"`
	FromBug          int                 `json:"fromBug"`
	Feedback         int                 `json:"feedback"`
	Title            string              `json:"title"`
	Keywords         string              `json:"keywords"`
	Type             string              `json:"type"`
	Category         string              `json:"category"`
	Pri              int                 `json:"pri"`
	Estimate         float64             `json:"estimate"`
	Status           string              `json:"status"`
	SubStatus        string              `json:"subStatus"`
	Color            string              `json:"color"`
	Stage            string              `json:"stage"`
	Mailto           []interface{}       `json:"mailto"`
	Lib              int                 `json:"lib"`
	FromStory        uint64              `json:"fromStory"`
	FromVersion      int                 `json:"fromVersion"`
	OpenedBy         *ZentaoAccount      `json:"openedBy"`
	OpenedDate       *helper.Iso8601Time `json:"openedDate"`
	AssignedTo       *ZentaoAccount      `json:"assignedTo"`
	AssignedDate     *helper.Iso8601Time `json:"assignedDate"`
	ApprovedDate     *helper.Iso8601Time `json:"approvedDate"`
	LastEditedBy     *ZentaoAccount      `json:"lastEditedBy"`
	LastEditedDate   *helper.Iso8601Time `json:"lastEditedDate"`
	ChangedBy        string              `json:"changedBy"`
	ChangedDate      *helper.Iso8601Time `json:"changedDate"`
	ReviewedBy       *ZentaoAccount      `json:"reviewedBy"`
	ReviewedDate     *helper.Iso8601Time `json:"reviewedDate"`
	ClosedBy         *ZentaoAccount      `json:"closedBy"`
	ClosedDate       *helper.Iso8601Time `json:"closedDate"`
	ClosedReason     string              `json:"closedReason"`
	ActivatedDate    *helper.Iso8601Time `json:"activatedDate"`
	ToBug            int                 `json:"toBug"`
	ChildStories     string              `json:"childStories"`
	LinkStories      string              `json:"linkStories"`
	LinkRequirements string              `json:"linkRequirements"`
	DuplicateStory   uint64              `json:"duplicateStory"`
	Version          int                 `json:"version"`
	StoryChanged     string              `json:"storyChanged"`
	FeedbackBy       string              `json:"feedbackBy"`
	NotifyEmail      string              `json:"notifyEmail"`
	URChanged        string              `json:"URChanged"`
	Deleted          bool                `json:"deleted"`
	PriOrder         string              `json:"priOrder"`
	PlanTitle        string              `json:"planTitle"`
	ProductStatus    string              `json:"productStatus"`
}

type ZentaoStory struct {
	common.NoPKModel
	ConnectionId uint64  `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ID           uint64  `json:"id" gorm:"primaryKey;type:BIGINT  NOT NULL" `
	Product      uint64  `json:"product"`
	Branch       int     `json:"branch"`
	Version      int     `json:"version"`
	OrderIn      int     `json:"order"`
	Vision       string  `json:"vision"`
	Parent       uint64  `json:"parent"`
	Module       int     `json:"module"`
	Plan         string  `json:"plan"`
	Source       string  `json:"source"`
	SourceNote   string  `json:"sourceNote"`
	FromBug      int     `json:"fromBug"`
	Feedback     int     `json:"feedback"`
	Title        string  `json:"title"`
	Keywords     string  `json:"keywords"`
	Type         string  `json:"type"`
	Category     string  `json:"category"`
	Pri          int     `json:"pri"`
	Estimate     float64 `json:"estimate"`
	Status       string  `json:"status"`
	SubStatus    string  `json:"subStatus"`
	Color        string  `json:"color"`
	Stage        string  `json:"stage"`
	//Mailto           []interface{} `json:"mailto"`
	Lib              int    `json:"lib"`
	FromStory        uint64 `json:"fromStory"`
	FromVersion      int    `json:"fromVersion"`
	OpenedById       uint64
	OpenedByName     string
	OpenedDate       *helper.Iso8601Time `json:"openedDate"`
	AssignedToId     uint64
	AssignedToName   string
	AssignedDate     *helper.Iso8601Time `json:"assignedDate"`
	ApprovedDate     *helper.Iso8601Time `json:"approvedDate"`
	LastEditedId     uint64
	LastEditedDate   *helper.Iso8601Time `json:"lastEditedDate"`
	ChangedDate      *helper.Iso8601Time `json:"changedDate"`
	ReviewedById     uint64              `json:"reviewedBy"`
	ReviewedDate     *helper.Iso8601Time `json:"reviewedDate"`
	ClosedId         uint64
	ClosedDate       *helper.Iso8601Time `json:"closedDate"`
	ClosedReason     string              `json:"closedReason"`
	ActivatedDate    *helper.Iso8601Time `json:"activatedDate"`
	ToBug            int                 `json:"toBug"`
	ChildStories     string              `json:"childStories"`
	LinkStories      string              `json:"linkStories"`
	LinkRequirements string              `json:"linkRequirements"`
	DuplicateStory   uint64              `json:"duplicateStory"`
	StoryChanged     string              `json:"storyChanged"`
	FeedbackBy       string              `json:"feedbackBy"`
	NotifyEmail      string              `json:"notifyEmail"`
	URChanged        string              `json:"URChanged"`
	Deleted          bool                `json:"deleted"`
	PriOrder         string              `json:"priOrder"`
	PlanTitle        string              `json:"planTitle"`
}

func (ZentaoStory) TableName() string {
	return "_tool_zentao_stories"
}
