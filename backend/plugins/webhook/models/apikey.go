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
	"gorm.io/datatypes"
	"time"
)

// ApiKey is the basic of api key management.
type ApiKey struct {
	ID           uint64         `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	Creator      string         `json:"creator"`
	CreatorEmail string         `json:"creatorEmail"`
	Updater      string         `json:"updater"`
	UpdaterEmail string         `json:"updater_email"`
	Name         string         `json:"name"`
	ApiKey       string         `json:"apiKey"`
	ExpiredAt    *time.Time     `json:"expiredAt"`
	AllowedPath  string         `json:"allowedPath"`
	Type         string         `json:"type"`
	Extra        datatypes.JSON `json:"extra"`
}

func (ApiKey) TableName() string {
	return "_devlake_api_keys"
}
