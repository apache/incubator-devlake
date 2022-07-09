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
	"time"
)

type BitbucketRepo struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	BitbucketId  string `gorm:"primaryKey"`
	Name         string `gorm:"type:varchar(255)"`
	HTMLUrl      string `gorm:"type:varchar(255)"`
	Description  string
	OwnerId      string     `json:"ownerId"`
	Language     string     `json:"language" gorm:"type:varchar(255)"`
	CreatedDate  time.Time  `json:"createdDate"`
	UpdatedDate  *time.Time `json:"updatedDate"`
	common.NoPKModel
}

func (BitbucketRepo) TableName() string {
	return "_tool_bitbucket_repos"
}
