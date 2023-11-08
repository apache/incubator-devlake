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
	"time"

	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type TapdIteration struct {
	ConnectionId uint64     `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Id           uint64     `gorm:"primaryKey;type:BIGINT  NOT NULL;autoIncrement:false" json:"id,string"`
	Name         string     `gorm:"type:varchar(255)" json:"name"`
	WorkspaceId  uint64     `json:"workspace_id,string"`
	Startdate    *time.Time `json:"startdate"`
	Enddate      *time.Time `json:"enddate"`
	Status       string     `gorm:"type:varchar(255)" json:"status"`
	ReleaseId    uint64     `gorm:"type:BIGINT" json:"release_id,string"`
	Description  string     `json:"description"`
	Creator      string     `gorm:"type:varchar(255)" json:"creator"`
	Created      *time.Time `json:"created"`
	Modified     *time.Time `json:"modified"`
	Completed    *time.Time `json:"completed"`
	Releaseowner string     `gorm:"type:varchar(255)" json:"releaseowner"`
	Launchdate   *time.Time `json:"launchdate"`
	Notice       string     `gorm:"type:varchar(255)" json:"notice"`
	Releasename  string     `gorm:"type:varchar(255)" json:"releasename"`
	archived.NoPKModel
}

type TapdWorkspaceIteration struct {
	archived.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey"`
	WorkspaceId  uint64 `gorm:"primaryKey"`
	IterationId  uint64 `gorm:"primaryKey"`
}

func (TapdIteration) TableName() string {
	return "_tool_tapd_iterations"
}

func (TapdWorkspaceIteration) TableName() string {
	return "_tool_tapd_workspace_iterations"
}
