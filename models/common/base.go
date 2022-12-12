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

package common

import (
	"regexp"
	"time"
)

type Model struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type NoPKModel struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	RawDataOrigin
}

// embedded fields for tool layer tables
type RawDataOrigin struct {
	// can be used for flushing outdated records from table
	RawDataParams string `gorm:"column:_raw_data_params;type:varchar(255);index" json:"_raw_data_params"`
	RawDataTable  string `gorm:"column:_raw_data_table;type:varchar(255)" json:"_raw_data_table"`
	// can be used for debugging
	RawDataId uint64 `gorm:"column:_raw_data_id" json:"_raw_data_id"`
	// we can store record index into this field, which is helpful for debugging
	RawDataRemark string `gorm:"column:_raw_data_remark" json:"_raw_data_remark"`
}

func NewNoPKModel() NoPKModel {
	now := time.Now()
	return NoPKModel{
		CreatedAt: now,
		UpdatedAt: now,
	}
}

var (
	DUPLICATE_REGEX = regexp.MustCompile(`(?i)\bduplicate\b`)
)

func IsDuplicateError(err error) bool {
	return err != nil && DUPLICATE_REGEX.MatchString(err.Error())
}
