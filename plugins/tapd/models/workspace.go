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

type TapdWorkspace struct {
	ConnectionId uint64         `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Id           uint64         `gorm:"primaryKey;type:BIGINT" json:"id,string"`
	Name         string         `gorm:"type:varchar(255)" json:"name"`
	PrettyName   string         `gorm:"type:varchar(255)" json:"pretty_name"`
	Category     string         `gorm:"type:varchar(255)" json:"category"`
	Status       string         `gorm:"type:varchar(255)" json:"status"`
	Description  string         `json:"description"`
	BeginDate    helper.CSTTime `json:"begin_date"`
	EndDate      helper.CSTTime `json:"end_date"`
	ExternalOn   string         `gorm:"type:varchar(255)" json:"external_on"`
	ParentId     uint64         `gorm:"type:BIGINT" json:"parent_id,string"`
	Creator      string         `gorm:"type:varchar(255)" json:"creator"`
	Created      helper.CSTTime `json:"created"`
	common.NoPKModel
}

func (TapdWorkspace) TableName() string {
	return "_tool_tapd_workspaces"
}
