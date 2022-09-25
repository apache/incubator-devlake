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

type ZentaoStories struct {
	common.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ID           int    `json:"id"gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Vision       string `json:"vision"`
	Parent       int    `json:"parent"`
	Product      int    `json:"product"`
	Branch       int    `json:"branch"`
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
	//Mailto           []interface{} `json:"mailto" gorm:"-:all"`
	Lib              int `json:"lib"`
	FromStory        int `json:"fromStory"`
	FromVersion      int `json:"fromVersion"`
	OpenedBy         `json:"openedBy"`
	OpenedDate       *helper.Iso8601Time `json:"openedDate"`
	AssignedTo       `json:"assignedTo"`
	AssignedDate     *helper.Iso8601Time `json:"assignedDate"`
	ApprovedDate     string              `json:"approvedDate"`
	LastEditedBy     `json:"lastEditedBy"`
	LastEditedDate   *helper.Iso8601Time `json:"lastEditedDate"`
	ChangedBy        string              `json:"changedBy"`
	ChangedDate      string              `json:"changedDate"`
	ReviewedBy       interface{}         `json:"reviewedBy" gorm:"-:all"`
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
	Version          int                 `json:"version"`
	StoryChanged     string              `json:"storyChanged"`
	FeedbackBy       string              `json:"feedbackBy"`
	NotifyEmail      string              `json:"notifyEmail"`
	URChanged        string              `json:"URChanged"`
	Deleted          bool                `json:"deleted"`
	Spec             string              `json:"spec"`
	Verify           string              `json:"verify"`
	Executions       Executions          `json:"executions" gorm:"-:all"`
	Tasks            []Tasks             `json:"tasks" gorm:"-:all"`
	//Stages           []interface{}       `json:"stages" gorm:"-:all"`
	PlanTitle []string `json:"planTitle" gorm:"-:all"`
	//Children         []interface{}       `json:"children" gorm:"-:all"`
	//Files            []interface{}       `json:"files" gorm:"-:all"`
	ProductName   string  `json:"productName"`
	ProductStatus string  `json:"productStatus"`
	ModuleTitle   string  `json:"moduleTitle"`
	Bugs          []Bugs  `json:"bugs" gorm:"-:all"`
	Cases         []Cases `json:"cases" gorm:"-:all"`
	//Requirements  []interface{} `json:"requirements" gorm:"-:all"`
	Actions    []Actions `json:"actions" gorm:"-:all"`
	PreAndNext `json:"preAndNext"`
}
type Executions struct {
	Num1 struct {
		Project int    `json:"project"`
		Name    string `json:"name"`
		Status  string `json:"status"`
		Type    string `json:"type"`
	} `json:"1"`
}
type Tasks struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Status     string `json:"status"`
	AssignedTo struct {
		ID       int    `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"assignedTo"`
}
type Bugs struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Pri      int    `json:"pri"`
	Severity int    `json:"severity"`
}
type Cases struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Pri    int    `json:"pri"`
	Status string `json:"status"`
}

type Actions struct {
	ID         int    `json:"id"`
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
	//History    []interface{} `json:"history"`
	Desc string `json:"desc"`
}
type AssignedTo struct {
	ID       int    `json:"id"`
	Account  string `json:"account"`
	Avatar   string `json:"avatar"`
	Realname string `json:"realname"`
}

type PreAndNext struct {
	Pre  string `json:"pre"`
	Next string `json:"next"`
}

func (ZentaoStories) TableName() string {
	return "_tool_zentao_stories"
}
