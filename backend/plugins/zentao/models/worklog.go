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

type ZentaoWorklog struct {
	ConnectionId uint64  `gorm:"primaryKey;type:BIGINT NOT NULL"`
	Id           int64   `gorm:"primaryKey;type:BIGINT NOT NULL;autoIncrement:false" json:"id"`
	ObjectId     int64   `json:"objectID"`
	ObjectType   string  `json:"objectType"`
	Project      int64   `json:"project"`
	Execution    int64   `json:"execution"`
	Product      string  `json:"product"`
	Account      string  `json:"account"`
	Work         string  `json:"work"`
	Vision       string  `json:"vision"`
	Date         string  `json:"date"`
	Left         float32 `json:"left"`
	Consumed     float32 `json:"consumed"`
	Begin        int64   `json:"begin"`
	End          int64   `json:"end"`
	Extra        *string `json:"extra"`
	Order        int64   `json:"order"`
	Deleted      string  `json:"deleted"`
	common.NoPKModel
}

func (ZentaoWorklog) TableName() string {
	return "_tool_zentao_worklogs"
}
