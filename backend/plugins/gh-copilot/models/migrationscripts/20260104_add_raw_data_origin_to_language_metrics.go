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

package migrationscripts

import (
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

// addRawDataOriginToCopilotLanguageMetrics ensures _tool_copilot_language_metrics includes RawDataOrigin columns.
// This is required by StatefulApiExtractor, which attaches provenance to extracted records.
type addRawDataOriginToCopilotLanguageMetrics struct{}

type ghCopilotLanguageMetrics20260104 struct {
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

	archived.NoPKModel
}

func (ghCopilotLanguageMetrics20260104) TableName() string {
	return "_tool_copilot_org_language_metrics"
}

func (script *addRawDataOriginToCopilotLanguageMetrics) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&ghCopilotLanguageMetrics20260104{},
	)
}

func (*addRawDataOriginToCopilotLanguageMetrics) Version() uint64 {
	return 20260104000000
}

func (*addRawDataOriginToCopilotLanguageMetrics) Name() string {
	return "copilot add raw data origin to language metrics"
}
