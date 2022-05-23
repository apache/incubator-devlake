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
)

type TestConnectionRequest struct {
	Endpoint string `json:"endpoint" validate:"required,url"`
	Auth     string `json:"auth" validate:"required"`
	Proxy    string `json:"proxy"`
}

type WorkspaceResponse struct {
	Id    uint64
	Title string
	Value string
}

type TapdConnection struct {
	common.Model
	Name             string `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	Endpoint         string `gorm:"type:varchar(255)" json:"endpoint"`
	BasicAuthEncoded string `gorm:"type:varchar(255)" json:"basicAuthEncoded"`
	RateLimit        int    `comment:"api request rate limt per hour" json:"rateLimit"`
}

type TapdConnectionDetail struct {
	TapdConnection
}

func (TapdConnection) TableName() string {
	return "_tool_tapd_connections"
}
