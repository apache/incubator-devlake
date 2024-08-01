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
	"github.com/apache/incubator-devlake/core/models/common"
	"time"
)

type AzuredevopsWorkItem struct {
	common.NoPKModel

	//AzuredevopsId int    `gorm:"primaryKey"`
	ConnectionId uint64 `gorm:"primaryKey"`
	WorkItemID   string `gorm:"primaryKey"`
	Title        string
	Type         string
	State        string
	CreatedDate  *time.Time
	ResolvedDate *time.Time
	ChangedDate  *time.Time
	CreatorName  string
	CreatorId    string
	AssigneeName string
	Area         string
	Url          string
	Severity     string
	Priority     string
	StoryPoint   float64
}

func (AzuredevopsWorkItem) TableName() string {

	return "_tool_azuredevops_go_workitem"
}

type AzuredevopsApiWorkItem struct {
	Id     int    `json:"id"`
	Rev    int    `json:"rev"`
	Url    string `json:"url"`
	Fields struct {
		SystemAreaPath     string              `json:"System.AreaPath"`
		SystemTeamProject  string              `json:"System.TeamProject"`
		SystemWorkItemType string              `json:"System.WorkItemType"`
		SystemState        string              `json:"System.State"`
		SystemReason       string              `json:"System.Reason"`
		SystemCreatedDate  *common.Iso8601Time `json:"System.CreatedDate"`
		SystemChangedDate  *common.Iso8601Time `json:"System.ChangedDate"`
		SystemTitle        string              `json:"System.Title"`
		SystemDescription  string              `json:"System.Description"`
		SystemAssignedTo   struct {
			DisplayName string `json:"displayName"`
			Id          string `json:"id"`
		} `json:"System.AssignedTo"`
		MicrosoftVSTSSchedulingEffort float64 `json:"Microsoft\.VSTS\.Scheduling\.Effort"`
		MicrosoftVSTSCommonPriority   string  `json:"Microsoft\.VSTS\.Common\.Priority"`
		MicrosoftVSTSCommonSeverity   string  `json:"Microsoft\.VSTS\.Common\.Severity"`
		SystemCreatedBy               struct {
			DisplayName string `json:"displayName"`
			Id          string `json:"id"`
		} `json:"System.CreatedBy"`
	} `json:"fields"`
}
