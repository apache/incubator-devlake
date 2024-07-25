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
)

type SubtaskState struct {
	Plugin  string `gorm:"primaryKey;type:varchar(255)" json:"plugin"`
	Subtask string `gorm:"primaryKey;type:varchar(255)" json:"subtask"`
	// Params is a json string to identitfy rows of a specific scope (jira board, github repo)
	Params string `gorm:"primaryKey;type:varchar(255);index" json:"params"`
	// PrevConfig stores the previous configuration of the subtask for determining should subtask run in Incremntal or FullSync mode
	PrevConfig string `json:"prevConfig"`
	// TimeAfter stores the previous timeAfter specified by the user for determining should subtask run in Incremntal or FullSync mode
	TimeAfter     *time.Time `json:"timeAfter"`
	PrevStartedAt *time.Time `json:"prevStartedAt"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

func (SubtaskState) TableName() string {
	return "_devlake_subtask_states"
}
