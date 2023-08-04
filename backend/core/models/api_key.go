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
	"time"
)

// ApiKey is the basic of api key management.
type ApiKey struct {
	common.Model
	common.Creator
	common.Updater
	Name        string     `json:"name"`
	ApiKey      string     `json:"apiKey"`
	ExpiredAt   *time.Time `json:"expiredAt"`
	AllowedPath string     `json:"allowedPath"`
	Type        string     `json:"type"`
	Extra       string     `json:"extra"`
}

func (ApiKey) TableName() string {
	return "_devlake_api_keys"
}

type ApiInputApiKey struct {
	Name        string     `json:"name" validate:"required,max=255"`
	Type        string     `json:"type" validate:"required"`
	AllowedPath string     `json:"allowedPath" validate:"required"`
	ExpiredAt   *time.Time `json:"expiredAt" `
}

type ApiOutputApiKey = ApiKey
