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

	"github.com/apache/incubator-devlake/core/models/common"
)

// CopilotLanguageMetrics represents engagement statistics broken down by editor and language.
type CopilotLanguageMetrics struct {
	ConnectionId uint64    `gorm:"primaryKey"`
	ScopeId      string    `gorm:"primaryKey;type:varchar(255)"`
	Date         time.Time `gorm:"primaryKey;type:date"`
	Editor       string    `gorm:"primaryKey;type:varchar(50)"`
	Language     string    `gorm:"primaryKey;type:varchar(50)"`

	EngagedUsers   int `json:"engagedUsers"`
	Suggestions    int `json:"suggestions"`
	Acceptances    int `json:"acceptances"`
	LinesSuggested int `json:"linesSuggested"`
	LinesAccepted  int `json:"linesAccepted"`
	common.RawDataOrigin
}

func (CopilotLanguageMetrics) TableName() string {
	return "_tool_copilot_language_metrics"
}
