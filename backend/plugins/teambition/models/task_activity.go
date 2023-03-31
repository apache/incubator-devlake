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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type TeambitionTaskActivity struct {
	ConnectionId      uint64           `gorm:"primaryKey;type:BIGINT"`
	Id                string           `gorm:"primaryKey;type:varchar(100)" json:"id"`
	ProjectId         string           `gorm:"primaryKey;type:varchar(100)" json:"projectId"`
	TaskId            string           `gorm:"primaryKey;type:varchar(100)" json:"taskId"`
	CreatorId         string           `gorm:"type:varchar(100)" json:"creatorId"`
	Action            string           `gorm:"type:varchar(100)" json:"action"`
	BoundToObjectId   string           `gorm:"type:varchar(100)" json:"boundToObjectId"`
	BoundToObjectType string           `gorm:"type:varchar(100)" json:"BoundToObjectType"`
	CreateTime        *api.Iso8601Time `json:"createTime"`
	UpdateTime        *api.Iso8601Time `json:"updateTime"`
	Content           string           `gorm:"type:text" json:"content"`

	common.NoPKModel
}

func (TeambitionTaskActivity) TableName() string {
	return "_tool_teambition_task_activities"
}
