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
	"time"
)

const (
	USER = "user"
)

type User struct {
	Name  string
	Email string
}

type Model struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Creator struct {
	Creator      string `json:"creator"`
	CreatorEmail string `json:"creatorEmail"`
}

type Updater struct {
	Updater      string `json:"updater"`
	UpdaterEmail string `json:"updaterEmail"`
}

// embedded fields for tool layer tables
type RawDataOrigin struct {
	// can be used for flushing outdated records from table
	RawDataParams string `gorm:"column:_raw_data_params;type:varchar(255);index" json:"_raw_data_params" mapstructure:"rawDataParams"`
	RawDataTable  string `gorm:"column:_raw_data_table;type:varchar(255)" json:"_raw_data_table" mapstructure:"rawDataTable"`
	// can be used for debugging
	RawDataId uint64 `gorm:"column:_raw_data_id" json:"_raw_data_id" mapstructure:"rawDataId"`
	// we can store record index into this field, which is helpful for debugging
	RawDataRemark string `gorm:"column:_raw_data_remark" json:"_raw_data_remark" mapstructure:"rawDataRemark"`
}

type GetRawDataOrigin interface {
	GetRawDataOrigin() *RawDataOrigin
}

func (c *RawDataOrigin) GetRawDataOrigin() *RawDataOrigin {
	return c
}

type NoPKModel struct {
	CreatedAt     time.Time `json:"createdAt" mapstructure:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" mapstructure:"updatedAt"`
	RawDataOrigin `swaggerignore:"true" mapstructure:",squash"`
}

func NewNoPKModel() NoPKModel {
	now := time.Now()
	return NoPKModel{
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type Scope struct {
	NoPKModel     `mapstructure:",squash"`
	ConnectionId  uint64 `json:"connectionId" gorm:"primaryKey" validate:"required" mapstructure:"connectionId,omitempty"`
	ScopeConfigId uint64 `json:"scopeConfigId,omitempty" mapstructure:"scopeConfigId,omitempty"`
}

// ScopeConnectionId implements plugin.ToolLayerScope.
func (s Scope) ScopeConnectionId() uint64 {
	return s.ConnectionId
}

func (s Scope) ScopeScopeConfigId() uint64 {
	return s.ScopeConfigId
}

type ScopeConfig struct {
	Model
	Entities     []string `gorm:"type:json;serializer:json" json:"entities" mapstructure:"entities"`
	ConnectionId uint64   `json:"connectionId" gorm:"index" validate:"required" mapstructure:"connectionId,omitempty"`
	Name         string   `mapstructure:"name" json:"name" gorm:"type:varchar(255)" validate:"required"`
}

func (s ScopeConfig) ScopeConfigConnectionId() uint64 {
	return s.ConnectionId
}
func (s ScopeConfig) ScopeConfigId() uint64 {
	return s.ID
}
