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

type JenkinsPipeline struct {
	common.NoPKModel
	// collected fields
	ConnectionId uint64 `gorm:"primaryKey"`
	DurationSec  uint64
	Name         string    `gorm:"type:varchar(255);primaryKey"`
	Result       string    // Result
	Status       string    // Result
	Timestamp    int64     // start time
	CreatedDate  time.Time // convered by timestamp
	CommitSha    string    `gorm:"primaryKey;type:varchar(255)"`
	Type         string    `gorm:"index;type:varchar(255)"`
	Building     bool
	Repo         string `gorm:"type:varchar(255);index"`
	FinishedDate *time.Time
}

func (JenkinsPipeline) TableName() string {
	return "_tool_jenkins_pipelines"
}
