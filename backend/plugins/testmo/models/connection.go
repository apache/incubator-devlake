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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

// TestmoConn holds the essential information to connect to the Testmo API
type TestmoConn struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.AccessToken    `mapstructure:",squash"`
}

func (connection TestmoConn) Sanitize() TestmoConn {
	connection.Token = ""
	return connection
}

type TestmoConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	TestmoConn            `mapstructure:",squash"`
}

func (TestmoConnection) TableName() string {
	return "_tool_testmo_connections"
}

func (connection TestmoConnection) Sanitize() TestmoConnection {
	connection.TestmoConn = connection.TestmoConn.Sanitize()
	return connection
}

func (connection *TestmoConnection) MergeFromRequest(target *TestmoConnection, body map[string]interface{}) error {
	token := target.Token
	if err := helper.DecodeMapStruct(body, target, true); err != nil {
		return err
	}
	modifiedToken := target.Token
	if modifiedToken == "" {
		target.Token = token
	}
	return nil
}

type TestmoResponse struct {
	Data   interface{} `json:"data"`
	Links  interface{} `json:"links"`
	Meta   interface{} `json:"meta"`
	Errors interface{} `json:"errors"`
}
