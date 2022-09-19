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

type ZentaoProject struct {
	common.NoPKModel
	ConnectionId  uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ID            int    `json:"id" gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Project       int    `json:"project"`
	Model         string `json:"model"`
	Type          string `json:"type"`
	Lifetime      string `json:"lifetime"`
	Budget        string `json:"budget"`
	BudgetUnit    string `json:"budgetUnit"`
	Attribute     string `json:"attribute"`
	Percent       int    `json:"percent"`
	Milestone     string `json:"milestone"`
	Output        string `json:"output"`
	Auth          string `json:"auth"`
	Parent        int    `json:"parent"`
	Path          string `json:"path"`
	Grade         int    `json:"grade"`
	Name          string `json:"name"`
	Code          string `json:"code"`
	Begin         string `json:"begin"`
	End           string `json:"end"`
	RealBegan     string `json:"realBegan"`
	RealEnd       string `json:"realEnd"`
	Days          int    `json:"days"`
	Status        string `json:"status"`
	SubStatus     string `json:"subStatus"`
	Pri           string `json:"pri"`
	Desc          string `json:"desc"`
	Version       int    `json:"version"`
	ParentVersion int    `json:"parentVersion"`
	PlanDuration  int    `json:"planDuration"`
	RealDuration  int    `json:"realDuration"`
	//OpenedBy       string    `json:"openedBy"`
	OpenedDate     time.Time  `json:"openedDate"`
	OpenedVersion  string     `json:"openedVersion"`
	LastEditedBy   string     `json:"lastEditedBy"`
	LastEditedDate *time.Time `json:"lastEditedDate,string"`
	ClosedBy       string     `json:"closedBy"`
	ClosedDate     *time.Time `json:"closedDate,string"`
	CanceledBy     string     `json:"canceledBy"`
	CanceledDate   *time.Time `json:"canceledDate,string"`
	SuspendedDate  string     `json:"suspendedDate"`
	PO             string     `json:"PO"`
	PM             `json:"PM"`
	QD             string `json:"QD"`
	RD             string `json:"RD"`
	Team           string `json:"team"`
	Acl            string `json:"acl"`
	Whitelist      `json:"whitelist" gorm:"-"`
	Order          int    `json:"order"`
	Vision         string `json:"vision"`
	DisplayCards   int    `json:"displayCards"`
	FluidBoard     string `json:"fluidBoard"`
	Deleted        bool   `json:"deleted"`
	Delay          int    `json:"delay"`
	Hours          `json:"hours"`
	TeamCount      int    `json:"teamCount"`
	LeftTasks      string `json:"leftTasks"`
	//TeamMembers   []interface{} `json:"teamMembers" gorm:"-"`
	TotalEstimate int `json:"totalEstimate"`
	TotalConsumed int `json:"totalConsumed"`
	TotalLeft     int `json:"totalLeft"`
	Progress      int `json:"progress"`
	TotalReal     int `json:"totalReal"`
}
type PM struct {
	PmId       int    `json:"id"`
	PmAccount  string `json:"account"`
	PmAvatar   string `json:"avatar"`
	PmRealname string `json:"realname"`
}
type Whitelist []struct {
	WhitelistID       int    `json:"id"`
	WhitelistAccount  string `json:"account"`
	WhitelistAvatar   string `json:"avatar"`
	WhitelistRealname string `json:"realname"`
}
type Hours struct {
	HoursTotalEstimate int `json:"totalEstimate"`
	HoursTotalConsumed int `json:"totalConsumed"`
	HoursTotalLeft     int `json:"totalLeft"`
	HoursProgress      int `json:"progress"`
	HoursTotalReal     int `json:"totalReal"`
}

func (ZentaoProject) TableName() string {
	return "_tool_zentao_projects"
}
