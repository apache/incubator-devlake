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

type SonarqubeMetrics struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	Id           string `gorm:"primaryKey"`
	MetricsKey   string `gorm:"type:varchar(100)"`
	Type         string `gorm:"type:varchar(100)"`
	Name         string `gorm:"type:varchar(100)"`
	Description  string
	Domain       string `gorm:"type:varchar(255)"`
	Direction    int
	Qualitative  bool
	Hidden       bool
	Custom       bool
	common.NoPKModel
}

func (SonarqubeMetrics) TableName() string {
	return "_tool_sonarqube_metrics"
}

type SonarqubeApiMetrics struct {
	Id          string `json:"id"`
	Key         string `json:"key"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Domain      string `json:"domain"`
	Direction   int    `json:"direction"`
	Qualitative bool   `json:"qualitative"`
	Hidden      bool   `json:"hidden"`
	Custom      bool   `json:"custom"`
}
