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
)

type ZentaoAccount struct {
	common.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ID           int64  `json:"id" gorm:"primaryKey;type:BIGINT  NOT NULL" `
	Account      string `json:"account" gorm:"type:varchar(100);index"`
	Avatar       string `json:"avatar" gorm:"type:varchar(255)"`
	Realname     string `json:"realname" gorm:"type:varchar(100);index"`
	Role         string `json:"role" gorm:"type:varchar(100);index"`
	Dept         int64  `json:"dept" gorm:"type:BIGINT  NOT NULL;index"`
}

func (ZentaoAccount) TableName() string {
	return "_tool_zentao_accounts"
}
