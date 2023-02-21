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

import "time"

type Model struct {
	ID        uint64     `gorm:"primaryKey" json:"id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
type BaseConnection struct {
	Name string `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	Model
}

// RestConnection implements the ApiConnection interface
type RestConnection struct {
	Endpoint         string `mapstructure:"endpoint" validate:"required" json:"endpoint"`
	Proxy            string `mapstructure:"proxy" json:"proxy"`
	RateLimitPerHour int    `comment:"api request rate limit per hour" json:"rateLimitPerHour"`
}

type BasicAuth struct {
	Username string `mapstructure:"username" validate:"required" json:"username"`
	Password string `mapstructure:"password" validate:"required" json:"password"`
}

type BambooConn struct {
	RestConnection `mapstructure:",squash"`
	//TODO you may need to use helper.BasicAuth instead of helper.AccessToken
	BasicAuth `mapstructure:",squash"`
}

// TODO Please modify the following code to fit your needs
// This object conforms to what the frontend currently sends.
type BambooConnection struct {
	BaseConnection `mapstructure:",squash"`
	BambooConn     `mapstructure:",squash"`
}

type TestConnectionRequest struct {
	Endpoint  string `json:"endpoint"`
	Proxy     string `json:"proxy"`
	BasicAuth `mapstructure:",squash"`
}

// This object conforms to what the frontend currently expects.
type BambooResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	BambooConnection
}

// Using User because it requires authentication.
type ApiResourcesResponse struct {
	Resources string `json:"resources" xml:"resources"`
}

func (BambooConnection) TableName() string {
	return "_tool_bamboo_connections"
}
