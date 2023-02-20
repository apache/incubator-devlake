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

type BambooPlan struct {
	ConnectionId              uint64  `gorm:"primaryKey"`
	PlanKey                   string  `json:"planKey" gorm:"primaryKey"`
	Name                      string  `json:"name"`
	Expand                    string  `json:"expand"`
	ProjectKey                string  `json:"projectKey" gorm:"index"`
	ProjectName               string  `json:"projectName"`
	Description               string  `json:"description"`
	ShortName                 string  `json:"shortName"`
	BuildName                 string  `json:"buildName"`
	ShortKey                  string  `json:"shortKey"`
	Type                      string  `json:"type"`
	Enabled                   bool    `json:"enabled"`
	Href                      string  `json:"href"`
	Rel                       string  `json:"rel"`
	IsFavourite               bool    `json:"isFavourite"`
	IsActive                  bool    `json:"isActive"`
	IsBuilding                bool    `json:"isBuilding"`
	AverageBuildTimeInSeconds float64 `json:"averageBuildTimeInSeconds"`
	archived.NoPKModel
}

func (b *BambooPlan) TableName() string {
	return "_tool_bamboo_plans"
}
