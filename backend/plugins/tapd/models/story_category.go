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
)

type TapdStoryCategory struct {
	ConnectionId uint64          `gorm:"primaryKey"`
	Id           int64           `gorm:"primaryKey;type:BIGINT NOT NULL;autoIncrement:false" json:"id,string"`
	Name         string          `json:"name" gorm:"type:varchar(255)"`
	Description  string          `json:"description"`
	ParentId     int64           `json:"parent_id,string"`
	Created      *common.CSTTime `json:"created"`
	Modified     *common.CSTTime `json:"modified"`
	common.NoPKModel
}

func (TapdStoryCategory) TableName() string {
	return "_tool_tapd_story_categories"
}
