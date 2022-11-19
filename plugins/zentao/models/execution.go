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

type ZentaoExecutionRes struct {
	ID            uint64     `json:"id"`
	Project       uint64     `json:"project"`
	Model         string     `json:"model"`
	Type          string     `json:"type"`
	Lifetime      string     `json:"lifetime"`
	Budget        string     `json:"budget"`
	BudgetUnit    string     `json:"budgetUnit"`
	Attribute     string     `json:"attribute"`
	Percent       int        `json:"percent"`
	Milestone     string     `json:"milestone"`
	Output        string     `json:"output"`
	Auth          string     `json:"auth"`
	Parent        uint64     `json:"parent"`
	Path          string     `json:"path"`
	Grade         int        `json:"grade"`
	Name          string     `json:"name"`
	Code          string     `json:"code"`
	PlanBegin     string     `json:"begin"`
	PlanEnd       string     `json:"end"`
	RealBegan     *time.Time `json:"realBegan"`
	RealEnd       *time.Time `json:"realEnd"`
	Days          int        `json:"days"`
	Status        string     `json:"status"`
	SubStatus     string     `json:"subStatus"`
	Pri           string     `json:"pri"`
	Description   string     `json:"desc"`
	Version       int        `json:"version"`
	ParentVersion int        `json:"parentVersion"`
	PlanDuration  int        `json:"planDuration"`
	RealDuration  int        `json:"realDuration"`
	OpenedBy      struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"openedBy"`
	OpenedDate    *time.Time `json:"openedDate"`
	OpenedVersion string     `json:"openedVersion"`
	LastEditedBy  struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"lastEditedBy"`
	LastEditedDate *time.Time `json:"lastEditedDate"`
	ClosedBy       struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"closedBy"`
	ClosedDate *time.Time `json:"closedDate"`
	CanceledBy struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"canceledBy"`
	CanceledDate  *time.Time `json:"canceledDate"`
	SuspendedDate string     `json:"suspendedDate"`
	PO            struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"PO"`
	PM struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"PM"`
	QD struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"QD"`
	RD struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"RD"`
	Team      string `json:"team"`
	Acl       string `json:"acl"`
	Whitelist []struct {
		ID       uint64 `json:"id"`
		Account  string `json:"account"`
		Avatar   string `json:"avatar"`
		Realname string `json:"realname"`
	} `json:"whitelist"`
	OrderIn       int    `json:"order"`
	Vision        string `json:"vision"`
	DisplayCards  int    `json:"displayCards"`
	FluidBoard    string `json:"fluidBoard"`
	Deleted       bool   `json:"deleted"`
	TotalHours    int    `json:"totalHours"`
	TotalEstimate int    `json:"totalEstimate"`
	TotalConsumed int    `json:"totalConsumed"`
	TotalLeft     int    `json:"totalLeft"`
	ProjectInfo   struct {
		ID             uint64 `json:"id"`
		Project        uint64 `json:"project"`
		Model          string `json:"model"`
		Type           string `json:"type"`
		Lifetime       string `json:"lifetime"`
		Budget         string `json:"budget"`
		BudgetUnit     string `json:"budgetUnit"`
		Attribute      string `json:"attribute"`
		Percent        int    `json:"percent"`
		Milestone      string `json:"milestone"`
		Output         string `json:"output"`
		Auth           string `json:"auth"`
		Parent         uint64 `json:"parent"`
		Path           string `json:"path"`
		Grade          int    `json:"grade"`
		Name           string `json:"name"`
		Code           string `json:"code"`
		PlanBegin      string `json:"begin"`
		PlanEnd        string `json:"end"`
		RealBegan      string `json:"realBegan"`
		RealEnd        string `json:"realEnd"`
		Days           int    `json:"days"`
		Status         string `json:"status"`
		SubStatus      string `json:"subStatus"`
		Pri            string `json:"pri"`
		Description    string `json:"desc"`
		Version        int    `json:"version"`
		ParentVersion  int    `json:"parentVersion"`
		PlanDuration   int    `json:"planDuration"`
		RealDuration   int    `json:"realDuration"`
		OpenedBy       string `json:"openedBy"`
		OpenedDate     string `json:"openedDate"`
		OpenedVersion  string `json:"openedVersion"`
		LastEditedBy   string `json:"lastEditedBy"`
		LastEditedDate string `json:"lastEditedDate"`
		ClosedBy       string `json:"closedBy"`
		ClosedDate     string `json:"closedDate"`
		CanceledBy     string `json:"canceledBy"`
		CanceledDate   string `json:"canceledDate"`
		SuspendedDate  string `json:"suspendedDate"`
		PO             string `json:"PO"`
		PM             string `json:"PM"`
		QD             string `json:"QD"`
		RD             string `json:"RD"`
		Team           string `json:"team"`
		Acl            string `json:"acl"`
		Whitelist      string `json:"whitelist"`
		OrderIn        int    `json:"order"`
		Vision         string `json:"vision"`
		DisplayCards   int    `json:"displayCards"`
		FluidBoard     string `json:"fluidBoard"`
		Deleted        string `json:"deleted"`
	} `json:"projectInfo"`
	Progress    int `json:"progress"`
	TeamMembers []struct {
		ID         uint64 `json:"id"`
		Root       int    `json:"root"`
		Type       string `json:"type"`
		Account    string `json:"account"`
		Role       string `json:"role"`
		Position   string `json:"position"`
		Limited    string `json:"limited"`
		Join       string `json:"join"`
		Days       int    `json:"days"`
		Hours      int    `json:"hours"`
		Estimate   string `json:"estimate"`
		Consumed   string `json:"consumed"`
		Left       string `json:"left"`
		OrderIn    int    `json:"order"`
		TotalHours int    `json:"totalHours"`
		UserID     uint64 `json:"userID"`
		Realname   string `json:"realname"`
	} `json:"teamMembers"`
	Products []struct {
		ID    uint64        `json:"id"`
		Name  string        `json:"name"`
		Plans []interface{} `json:"plans"`
	} `json:"products"`
	CaseReview bool `json:"caseReview"`
}

type ZentaoExecution struct {
	ConnectionId   uint64     `gorm:"primaryKey"`
	Id             uint64     `json:"id" gorm:"primaryKey"`
	Project        uint64     `json:"project"`
	Model          string     `json:"model"`
	Type           string     `json:"type"`
	Lifetime       string     `json:"lifetime"`
	Budget         string     `json:"budget"`
	BudgetUnit     string     `json:"budgetUnit"`
	Attribute      string     `json:"attribute"`
	Percent        int        `json:"percent"`
	Milestone      string     `json:"milestone"`
	Output         string     `json:"output"`
	Auth           string     `json:"auth"`
	Parent         uint64     `json:"parent"`
	Path           string     `json:"path"`
	Grade          int        `json:"grade"`
	Name           string     `json:"name"`
	Code           string     `json:"code"`
	PlanBegin      string     `json:"begin"`
	PlanEnd        string     `json:"end"`
	RealBegan      *time.Time `json:"realBegan"`
	RealEnd        *time.Time `json:"realEnd"`
	Days           int        `json:"days"`
	Status         string     `json:"status"`
	SubStatus      string     `json:"subStatus"`
	Pri            string     `json:"pri"`
	Description    string     `json:"desc"`
	Version        int        `json:"version"`
	ParentVersion  int        `json:"parentVersion"`
	PlanDuration   int        `json:"planDuration"`
	RealDuration   int        `json:"realDuration"`
	OpenedById     uint64
	OpenedDate     *time.Time `json:"openedDate"`
	OpenedVersion  string     `json:"openedVersion"`
	LastEditedById uint64
	LastEditedDate *time.Time `json:"lastEditedDate"`
	ClosedById     uint64
	ClosedDate     *time.Time `json:"closedDate"`
	CanceledById   uint64
	CanceledDate   *time.Time `json:"canceledDate"`
	SuspendedDate  string     `json:"suspendedDate"`
	POId           uint64
	PMId           uint64
	QDId           uint64
	RDId           uint64
	Team           string `json:"team"`
	Acl            string `json:"acl"`
	OrderIn        int    `json:"order"`
	Vision         string `json:"vision"`
	DisplayCards   int    `json:"displayCards"`
	FluidBoard     string `json:"fluidBoard"`
	Deleted        bool   `json:"deleted"`
	TotalHours     int    `json:"totalHours"`
	TotalEstimate  int    `json:"totalEstimate"`
	TotalConsumed  int    `json:"totalConsumed"`
	TotalLeft      int    `json:"totalLeft"`
	ProjectId      uint64
	Progress       int  `json:"progress"`
	CaseReview     bool `json:"caseReview"`
	common.NoPKModel
}

func (ZentaoExecution) TableName() string {
	return "_tool_zentao_executions"
}
