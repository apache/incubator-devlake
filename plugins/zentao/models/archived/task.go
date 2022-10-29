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

type ZentaoTask struct {
	archived.NoPKModel
	ConnectionId  uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ExecutionId   uint64 `json:"execution_id"`
	ID            int    `json:"id" gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Project       int    `json:"project"`
	Parent        int    `json:"parent"`
	Execution     int    `json:"execution"`
	Module        int    `json:"module"`
	Design        int    `json:"design"`
	Story         int    `json:"story"`
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
	Left          int    `json:"left"`
	Deadline      string `json:"deadline"`
	Status        string `json:"status"`
	SubStatus     string `json:"subStatus"`
	Color         string `json:"color"`
	//Mailto        interface{} `json:"mailto"`
	Desc               string `json:"desc"`
	Version            int    `json:"version"`
	OpenedBy           `json:"openedBy"`
	OpenedDate         *helper.Iso8601Time `json:"openedDate"`
	AssignedTo         `json:"assignedTo"`
	AssignedDate       *helper.Iso8601Time `json:"assignedDate"`
	EstStarted         string              `json:"estStarted"`
	RealStarted        *helper.Iso8601Time `json:"realStarted"`
	FinishedBy         `json:"finishedBy"`
	FinishedDate       *helper.Iso8601Time `json:"finishedDate"`
	FinishedList       string              `json:"finishedList"`
	CanceledBy         `json:"canceledBy"`
	CanceledDate       *helper.Iso8601Time `json:"canceledDate"`
	ClosedBy           *helper.Iso8601Time `json:"closedBy"`
	ClosedDate         *helper.Iso8601Time `json:"closedDate"`
	PlanDuration       int                 `json:"planDuration"`
	RealDuration       int                 `json:"realDuration"`
	ClosedReason       string              `json:"closedReason"`
	LastEditedBy       `json:"lastEditedBy"`
	LastEditedDate     *helper.Iso8601Time `json:"lastEditedDate"`
	ActivatedDate      string              `json:"activatedDate"`
	Order              int                 `json:"order"`
	Repo               int                 `json:"repo"`
	Mr                 int                 `json:"mr"`
	Entry              string              `json:"entry"`
	Lines              string              `json:"lines"`
	V1                 string              `json:"v1"`
	V2                 string              `json:"v2"`
	Deleted            bool                `json:"deleted"`
	Vision             string              `json:"vision"`
	StoryID            int                 `json:"storyID"`
	StoryTitle         string              `json:"storyTitle"`
	Product            int                 `json:"product"`
	Branch             int                 `json:"branch"`
	LatestStoryVersion int                 `json:"latestStoryVersion"`
	StoryStatus        string              `json:"storyStatus"`
	AssignedToRealName string              `json:"assignedToRealName"`
	PriOrder           string              `json:"priOrder"`
	NeedConfirm        bool                `json:"needConfirm"`
	ProductType        string              `json:"productType"`
	Progress           int                 `json:"progress"`
}

type FinishedBy struct {
	FinishedByID       int    `json:"id"`
	FinishedByAccount  string `json:"account"`
	FinishedByAvatar   string `json:"avatar"`
	FinishedByRealname string `json:"realname"`
}

func (ZentaoTask) TableName() string {
	return "_tool_zentao_tasks"
}
