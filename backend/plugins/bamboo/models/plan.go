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
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.ToolLayerScope = (*BambooPlan)(nil)
var _ plugin.ApiScope = (*ApiBambooPlan)(nil)

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
	ScopeConfigId             uint64  `json:"scopeConfigId"`
	common.NoPKModel
}

func (p BambooPlan) ScopeId() string {
	return p.PlanKey
}

func (p BambooPlan) ScopeName() string {
	return p.Name
}

func (p BambooPlan) ScopeFullName() string {
	return p.Name
}

func (p BambooPlan) ScopeParams() interface{} {
	return &BambooApiParams{
		ConnectionId: p.ConnectionId,
		PlanKey:      p.PlanKey,
	}
}

func (p BambooPlan) TableName() string {
	return "_tool_bamboo_plans"
}

type ApiBambooPlan struct {
	ShortName string `json:"shortName"`
	ShortKey  string `json:"shortKey"`
	Type      string `json:"type"`
	Enabled   bool   `json:"enabled"`
	Link      struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"link"`
	Key     string `json:"key"`
	Name    string `json:"name"`
	PlanKey struct {
		Key string `json:"key"`
	} `json:"planKey"`
}

//type ApiBambooPlan struct {
//	Expand                    string `json:"expand"`
//	Description               string `json:"description"`
//	ShortName                 string `json:"shortName"`
//	BuildName                 string `json:"buildName"`
//	ShortKey                  string `json:"shortKey"`
//	Type                      string `json:"type"`
//	Enabled                   bool   `json:"enabled"`
//	ProjectKey                string `json:"projectKey"`
//	ProjectName               string `json:"projectName"`
//	ApiBambooLink             `json:"link"`
//	IsFavourite               bool    `json:"isFavourite"`
//	IsActive                  bool    `json:"isActive"`
//	IsBuilding                bool    `json:"isBuilding"`
//	AverageBuildTimeInSeconds float64 `json:"averageBuildTimeInSeconds"`
//	Key                       string  `json:"key"`
//	Name                      string  `json:"name"`
//}

func (p ApiBambooPlan) ConvertApiScope() plugin.ToolLayerScope {
	return &BambooPlan{
		PlanKey:   p.Key,
		Name:      p.Name,
		ShortName: p.ShortName,
		ShortKey:  p.ShortKey,
		Type:      p.Type,
		Enabled:   p.Enabled,
		Href:      p.Link.Href,
		Rel:       p.Link.Rel,
	}
}

type ApiSearchResult struct {
	Id         string        `json:"id"`
	EntityType string        `json:"entityType"`
	Entity     ApiBambooPlan `json:"entity"`
}

type ApiBambooSearchPlanResponse struct {
	ApiBambooSizeData `json:"squash"`
	SearchResults     []ApiSearchResult `json:"searchResults"`
}
