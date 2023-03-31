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

type TeambitionAccount struct {
	ConnectionId   uint64           `gorm:"primaryKey;type:BIGINT"`
	UserId         string           `gorm:"primaryKey;type:varchar(100)" json:"userId"`
	MemberId       string           `gorm:"type:varchar(100)" json:"memberId"`
	IsDisabled     int              `json:"isDisabled"`
	Role           uint64           `json:"role"`
	AvatarUrl      string           `gorm:"type:varchar(255)" json:"avatarUrl"`
	Birthday       string           `gorm:"type:varchar(100)" json:"birthday"`
	City           string           `gorm:"type:varchar(100)" json:"city"`
	Province       string           `gorm:"type:varchar(100)" json:"province"`
	Country        string           `gorm:"type:varchar(100)" json:"country"`
	Email          string           `gorm:"type:varchar(255)" json:"email"`
	EntryTime      *api.Iso8601Time `json:"entryTime"`
	Name           string           `gorm:"index;type:varchar(255)" json:"name"`
	Phone          string           `gorm:"type:varchar(100)" json:"phone"`
	Title          string           `gorm:"type:varchar(255)" json:"title"`
	Pinyin         string           `gorm:"type:varchar(255)" json:"pinyin"`
	Py             string           `gorm:"type:varchar(255)" json:"py"`
	StaffType      string           `gorm:"type:varchar(255)" json:"staffType"`
	EmployeeNumber string           `gorm:"type:varchar(100)" json:"employeeNumber"`

	common.NoPKModel
}

func (TeambitionAccount) TableName() string {
	return "_tool_teambition_accounts"
}
