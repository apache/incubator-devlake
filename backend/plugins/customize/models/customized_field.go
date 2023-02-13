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
	"time"

	"github.com/apache/incubator-devlake/core/dal"
)

type CustomizedField struct {
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	TbName      string         `gorm:"primaryKey;type:varchar(255)"` // avoid conflicting with the method `TableName()`
	ColumnName  string         `gorm:"primaryKey;type:varchar(255)"`
	DisplayName string         `gorm:"type:varchar(255)"`
	DataType    dal.ColumnType `gorm:"type:varchar(255)"`
	Description string
}

func (t *CustomizedField) TableName() string {
	return "_tool_customized_fields"
}
