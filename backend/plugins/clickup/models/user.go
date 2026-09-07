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
	"github.com/apache/devlake/core/models/common"
)

// ClickUpUser is a ClickUp user (tool layer), converted to crossdomain.Account.
type ClickUpUser struct {
	ConnectionId   uint64 `gorm:"primaryKey"`
	Id             string `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Username       string `gorm:"type:varchar(255)" json:"username"`
	Email          string `gorm:"type:varchar(255)" json:"email"`
	Color          string `gorm:"type:varchar(50)" json:"color"`
	ProfilePicture string `gorm:"type:varchar(255)" json:"profilePicture"`
	common.NoPKModel
}

func (ClickUpUser) TableName() string {
	return "_tool_clickup_users"
}
