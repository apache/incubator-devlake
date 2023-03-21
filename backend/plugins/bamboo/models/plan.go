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
	common.NoPKModel
}

func (apiRes *ApiBambooPlan) Convert() *BambooPlan {
	return &BambooPlan{
		PlanKey:                   apiRes.Key,
		Name:                      apiRes.Name,
		Expand:                    apiRes.Expand,
		ProjectKey:                apiRes.ProjectKey,
		ProjectName:               apiRes.ProjectName,
		Description:               apiRes.Description,
		ShortName:                 apiRes.ShortName,
		BuildName:                 apiRes.BuildName,
		ShortKey:                  apiRes.ShortKey,
		Type:                      apiRes.Type,
		Enabled:                   apiRes.Enabled,
		Href:                      apiRes.Href,
		Rel:                       apiRes.Rel,
		IsFavourite:               apiRes.IsFavourite,
		IsActive:                  apiRes.IsActive,
		IsBuilding:                apiRes.IsBuilding,
		AverageBuildTimeInSeconds: apiRes.AverageBuildTimeInSeconds,
	}
}

func (BambooPlan) TableName() string {
	return "_tool_bamboo_plans"
}

type ApiBambooPlan struct {
	Expand                    string `json:"expand"`
	Description               string `json:"description"`
	ShortName                 string `json:"shortName"`
	BuildName                 string `json:"buildName"`
	ShortKey                  string `json:"shortKey"`
	Type                      string `json:"type"`
	Enabled                   bool   `json:"enabled"`
	ProjectKey                string `json:"projectKey"`
	ProjectName               string `json:"projectName"`
	ApiBambooLink             `json:"link"`
	IsFavourite               bool    `json:"isFavourite"`
	IsActive                  bool    `json:"isActive"`
	IsBuilding                bool    `json:"isBuilding"`
	AverageBuildTimeInSeconds float64 `json:"averageBuildTimeInSeconds"`
	Key                       string  `json:"key"`
	Name                      string  `json:"name"`
}
