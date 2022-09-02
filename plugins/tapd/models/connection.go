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
	"encoding/base64"
	"fmt"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type TestConnectionRequest struct {
	Endpoint         string `json:"endpoint"`
	Proxy            string `json:"proxy"`
	helper.BasicAuth `mapstructure:",squash"`
}

type WorkspaceResponse struct {
	Id    uint64
	Title string
	Value string
}

type BasicAuth struct {
	Username string `mapstructure:"username" validate:"required" json:"username"`
	Password string `mapstructure:"password" validate:"required" json:"password"`
}

func (ba BasicAuth) GetEncodedToken() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", ba.Username, ba.Password)))
}

type TapdConnection struct {
	helper.RestConnection `mapstructure:",squash"`
	BasicAuth             `mapstructure:",squash"`
}

type TapdConnectionDetail struct {
	TapdConnection
}

func (TapdConnection) TableName() string {
	return "_tool_tapd_connections"
}
