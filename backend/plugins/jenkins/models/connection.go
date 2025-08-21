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

// JenkinsConn holds the essential information to connect to the Jenkins API
type JenkinsConn struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.BasicAuth      `mapstructure:",squash"`
}

func (connection JenkinsConn) Sanitize() JenkinsConn {
	connection.Password = ""
	return connection
}

// JenkinsConnection holds JenkinsConn plus ID/Name for database storage
type JenkinsConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	JenkinsConn           `mapstructure:",squash"`
}

func (JenkinsConnection) TableName() string {
	return "_tool_jenkins_connections"
}

func (connection JenkinsConnection) Sanitize() JenkinsConnection {
	connection.JenkinsConn = connection.JenkinsConn.Sanitize()
	return connection
}

func (connection *JenkinsConnection) MergeFromRequest(target *JenkinsConnection, body map[string]interface{}) error {
	password := target.Password
	if err := helper.DecodeMapStruct(body, target, true); err != nil {
		return err
	}
	modifiedPassword := target.Password
	if modifiedPassword == "" {
		target.Password = password
	}
	return nil
}
