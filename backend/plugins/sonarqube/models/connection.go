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

import helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"

// This object conforms to what the frontend currently sends.
type SonarqubeConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	helper.RestConnection `mapstructure:",squash"`
	// For sonarqube, we can `use user_token:`
	helper.AccessToken `mapstructure:",squash"`
}

type TestConnectionRequest struct {
	Endpoint           string `json:"endpoint"`
	Proxy              string `json:"proxy"`
	helper.AccessToken `mapstructure:",squash"`
}

// This object conforms to what the frontend currently expects.
type SonarqubeResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	SonarqubeConnection
}

// Using User because it requires authentication.
type ApiUserResponse struct {
	Id   int
	Name string `json:"name"`
}

func (SonarqubeConnection) TableName() string {
	return "_tool_sonarqube_connections"
}
