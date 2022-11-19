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

type ZentaoTaskRes struct {
	Id            uint64 `json:"id"`
	Project       uint64 `json:"project"`
	Parent        uint64 `json:"parent"`
	Execution     uint64 `json:"execution"`
	Module        int    `json:"module"`
	Design        int    `json:"design"`
	Story         uint64 `json:"story"`
	StoryVersion  int    `json:"storyVersion"`
	DesignVersion int    `json:"designVersion"`
	FromBug       int    `json:"fromBug"`
	Feedback      int    `json:"feedback"`
	FromIssue     int    `json:"fromIssue"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Mode          string `json:"mode"`
	Pri           int    `json:"pri"`
	Estimate      int    `json:"estimate"`
	Consumed      int    `json:"consumed"`
	Deadline      string `json:"deadline"`
	Status        string `json:"status"`
	SubStatus     string `json:"subStatus"`
	Color         string `json:"color"`
	Mailto        []struct {
		Id       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"mailto"`
	Description string `json:"desc"`
	Version     int    `json:"version"`
	OpenedBy    struct {
		Id       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"openedBy"`
	OpenedDate *time.Time `json:"openedDate"`
	AssignedTo struct {
		Id       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"assignedTo"`
	AssignedDate *time.Time `json:"assignedDate"`
	EstStarted   string     `json:"estStarted"`
	RealStarted  *time.Time `json:"realStarted"`
	FinishedBy   struct {
		Id       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"finishedBy"`
	FinishedDate *time.Time `json:"finishedDate"`
	FinishedList string     `json:"finishedList"`
	CanceledBy   struct {
		Id       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"canceledBy"`
	CanceledDate *time.Time `json:"canceledDate"`
	ClosedBy     struct {
		Id       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"closedBy"`
	ClosedDate   *time.Time `json:"closedDate"`
	PlanDuration int        `json:"planDuration"`
	RealDuration int        `json:"realDuration"`
	ClosedReason string     `json:"closedReason"`
	LastEditedBy struct {
		Id       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"lastEditedBy"`
	LastEditedDate *time.Time `json:"lastEditedDate"`
	ActivatedDate  string     `json:"activatedDate"`
	OrderIn        int        `json:"order"`
	Repo           int        `json:"repo"`
	Mr             int        `json:"mr"`
	Entry          string     `json:"entry"`
	Lines          string     `json:"lines"`
	V1             string     `json:"v1"`
	V2             string     `json:"v2"`
	Deleted        bool       `json:"deleted"`
	Vision         string     `json:"vision"`
	StoryID        uint64     `json:"storyID"`
	StoryTitle     string     `json:"storyTitle"`
	Branch         interface {
	} `json:"branch"`
	LatestStoryVersion interface {
	} `json:"latestStoryVersion"`
	StoryStatus interface {
	} `json:"storyStatus"`
	AssignedToRealName string `json:"assignedToRealName"`
	PriOrder           string `json:"priOrder"`
	Delay              int    `json:"delay"`
	NeedConfirm        bool   `json:"needConfirm"`
	Progress           int    `json:"progress"`
}

type ZentaoTask struct {
	common.NoPKModel
	ConnectionId  uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ExecutionId   uint64 `json:"execution_id"`
	ID            uint64 `json:"id" gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Project       uint64 `json:"project"`
	Parent        uint64 `json:"parent"`
	Execution     uint64 `json:"execution"`
	Module        int    `json:"module"`
	Design        int    `json:"design"`
	Story         uint64 `json:"story"`
	StoryVersion  int    `json:"storyVersion"`
	DesignVersion int    `json:"designVersion"`
	FromBug       int    `json:"fromBug"`
	Feedback      int    `json:"feedback"`
	FromIssue     int    `json:"fromIssue"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Mode          string `json:"mode"`
	Pri           int    `json:"pri"`
	Estimate      int    `json:"estimate"`
	Consumed      int    `json:"consumed"`
	Deadline      string `json:"deadline"`
	Status        string `json:"status"`
	SubStatus     string `json:"subStatus"`
	Color         string `json:"color"`
	//Mailto        interface{} `json:"mailto"`
	Description        string `json:"desc"`
	Version            int    `json:"version"`
	OpenedById         uint64
	OpenedByName       string
	OpenedDate         *time.Time `json:"openedDate"`
	AssignedToId       uint64
	AssignedToName     string
	AssignedDate       *time.Time `json:"assignedDate"`
	EstStarted         string     `json:"estStarted"`
	RealStarted        *time.Time `json:"realStarted"`
	FinishedId         uint64
	FinishedDate       *time.Time `json:"finishedDate"`
	FinishedList       string     `json:"finishedList"`
	CanceledId         uint64
	CanceledDate       *time.Time `json:"canceledDate"`
	ClosedById         uint64
	ClosedDate         *time.Time `json:"closedDate"`
	PlanDuration       int        `json:"planDuration"`
	RealDuration       int        `json:"realDuration"`
	ClosedReason       string     `json:"closedReason"`
	LastEditedId       uint64
	LastEditedDate     *time.Time `json:"lastEditedDate"`
	ActivatedDate      string     `json:"activatedDate"`
	OrderIn            int        `json:"order"`
	Repo               int        `json:"repo"`
	Mr                 int        `json:"mr"`
	Entry              string     `json:"entry"`
	Lines              string     `json:"lines"`
	V1                 string     `json:"v1"`
	V2                 string     `json:"v2"`
	Deleted            bool       `json:"deleted"`
	Vision             string     `json:"vision"`
	StoryID            uint64     `json:"storyID"`
	StoryTitle         string     `json:"storyTitle"`
	Branch             int        `json:"branch"`
	LatestStoryVersion int        `json:"latestStoryVersion"`
	StoryStatus        string     `json:"storyStatus"`
	AssignedToRealName string     `json:"assignedToRealName"`
	PriOrder           string     `json:"priOrder"`
	NeedConfirm        bool       `json:"needConfirm"`
	Progress           int        `json:"progress"`
}

func (ZentaoTask) TableName() string {
	return "_tool_zentao_tasks"
}
