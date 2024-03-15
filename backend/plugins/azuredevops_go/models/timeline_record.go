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

type AzuredevopsTimelineRecord struct {
	common.NoPKModel

	ConnectionId uint64 `gorm:"primaryKey"`
	RecordId     string `gorm:"primaryKey"`
	BuildId      int    `gorm:"primaryKey"`
	ParentId     string
	Type         string
	Name         string
	StartTime    *time.Time
	FinishTime   *time.Time
	State        string
	Result       string
	ChangeId     int
	LastModified string
}

func (AzuredevopsTimelineRecord) TableName() string {
	return "_tool_azuredevops_go_timeline_records"
}

type AzuredevopsApiTimelineRecord struct {
	Id           string     `json:"id"`
	ParentId     string     `json:"parentId"`
	Type         string     `json:"type"`
	Name         string     `json:"name"`
	StartTime    *time.Time `json:"startTime"`
	FinishTime   *time.Time `json:"finishTime"`
	State        string     `json:"state"`
	Result       string     `json:"result"`
	ResultCode   string     `json:"resultCode"`
	ChangeId     int        `json:"changeId"`
	LastModified string     `json:"lastModified"`
	Identifier   string     `json:"identifier"`
}
