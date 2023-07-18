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
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type TeambitionAccount struct {
	ConnectionId   uint64 `gorm:"primaryKey;type:BIGINT"`
	UserId         string `gorm:"primaryKey;type:varchar(100)" json:"userId"`
	MemberId       string `gorm:"type:varchar(100)"`
	IsDisabled     int    `json:"isDisabled"`
	Role           uint64
	AvatarUrl      string           `gorm:"type:varchar(255)"`
	Birthday       string           `gorm:"type:varchar(100)"`
	City           string           `gorm:"type:varchar(100)"`
	Province       string           `gorm:"type:varchar(100)"`
	Country        string           `gorm:"type:varchar(100)"`
	Email          string           `gorm:"type:varchar(255)"`
	EntryTime      *api.Iso8601Time `json:"entryTime"`
	Name           string           `gorm:"index;type:varchar(255)"`
	Phone          string           `gorm:"type:varchar(100)"`
	Title          string           `gorm:"type:varchar(255)"`
	Pinyin         string           `gorm:"type:varchar(255)"`
	Py             string           `gorm:"type:varchar(255)"`
	StaffType      string           `gorm:"type:varchar(255)"`
	EmployeeNumber string           `gorm:"type:varchar(100)"`

	archived.NoPKModel
}

func (TeambitionAccount) TableName() string {
	return "_tool_teambition_accounts"
}
