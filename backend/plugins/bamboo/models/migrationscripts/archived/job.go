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

package archived

import (
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type BambooJob struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	JobKey       string `gorm:"primaryKey"`
	Id           string
	Name         string `json:"name"`
	PlanKey      string `json:"planKey"`
	PlanName     string `json:"planName"`
	ProjectKey   string `gorm:"index"`
	ProjectName  string `json:"projectName"`
	Description  string `json:"description"`
	BranchName   string `json:"branchName"`
	StageName    string `json:"stageName"`
	Type         string `json:"type"`
	archived.NoPKModel
}

func (BambooJob) TableName() string {
	return "_tool_bamboo_jobs"
}
