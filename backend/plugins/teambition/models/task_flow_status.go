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

type TeambitionTaskFlowStatus struct {
	ConnectionId                uint64           `gorm:"primaryKey;type:BIGINT"`
	ProjectId                   string           `gorm:"primaryKey;type:varchar(100)" json:"projectId"`
	Id                          string           `gorm:"primaryKey;type:varchar(100)" json:"id"`
	Name                        string           `gorm:"type:varchar(100)" json:"name"`
	Pos                         int64            `json:"pos"`
	TaskflowId                  string           `gorm:"type:varchar(100)" json:"taskflowId"`
	RejectStatusIds             []string         `gorm:"serializer:json;type:text" json:"rejectStatusIds"`
	Kind                        string           `gorm:"varchar(100)" json:"kind"`
	CreatorId                   string           `gorm:"varchar(100)" json:"creatorId"`
	IsDeleted                   bool             `json:"isDeleted"`
	IsTaskflowstatusruleexector bool             `json:"isTaskflowstatusruleexector"`
	Created                     *api.Iso8601Time `json:"created"`
	Updated                     *api.Iso8601Time `json:"updated"`

	common.NoPKModel
}

func (TeambitionTaskFlowStatus) TableName() string {
	return "_tool_teambition_task_flow_status"
}
