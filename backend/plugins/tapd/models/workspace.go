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
	"fmt"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
)

type WorkspacesResponse struct {
	Status int `json:"status"`
	Data   []struct {
		TapdWorkspace `json:"Workspace"`
	} `json:"data"`
	Info string `json:"info"`
}

type WorkspaceResponse struct {
	Status int `json:"status"`
	Data   struct {
		TapdWorkspace `json:"Workspace"`
	} `json:"data"`
	Info string `json:"info"`
}

type TapdWorkspace struct {
	common.Scope
	Id          uint64          `gorm:"primaryKey;type:BIGINT" mapstructure:"id" json:"id"`
	Name        string          `gorm:"type:varchar(255)" mapstructure:"name" json:"name"`
	PrettyName  string          `gorm:"type:varchar(255)" mapstructure:"pretty_name" json:"pretty_name"`
	Category    string          `gorm:"type:varchar(255)" mapstructure:"category" json:"category"`
	Status      string          `gorm:"type:varchar(255)" mapstructure:"status" json:"status"`
	Description string          `mapstructure:"description" json:"description"`
	BeginDate   *common.CSTTime `mapstructure:"begin_date" json:"begin_date"`
	EndDate     *common.CSTTime `mapstructure:"end_date" json:"end_date"`
	ExternalOn  string          `gorm:"type:varchar(255)" mapstructure:"external_on" json:"external_on"`
	ParentId    uint64          `gorm:"type:BIGINT" mapstructure:"parent_id,string" json:"parent_id"`
	Creator     string          `gorm:"type:varchar(255)" mapstructure:"creator" json:"creator"`
	Created     *common.CSTTime `mapstructure:"created" json:"created"`
}

func (TapdWorkspace) TableName() string {
	return "_tool_tapd_workspaces"
}

func (w TapdWorkspace) ScopeId() string {
	return fmt.Sprintf(`%d`, w.Id)
}

func (w TapdWorkspace) ScopeName() string {
	return w.Name
}

func (w TapdWorkspace) ScopeFullName() string {
	return w.Name
}

func (w TapdWorkspace) ScopeParams() interface{} {
	return &TapdApiParams{
		ConnectionId: w.ConnectionId,
		WorkspaceId:  w.Id,
	}
}

func (w TapdWorkspace) ConvertApiScope() plugin.ToolLayerScope {
	return w
}

type TapdApiParams struct {
	ConnectionId uint64
	WorkspaceId  uint64
}
