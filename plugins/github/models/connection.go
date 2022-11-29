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
	"github.com/apache/incubator-devlake/plugins/helper"
)

type TestConnectionRequest struct {
	Endpoint string `json:"endpoint" validate:"required,url"`
	Token    string `json:"token" validate:"required"`
	Proxy    string `json:"proxy"`
}

type GithubConnection struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.AccessToken    `mapstructure:",squash"`
	EnableGraphql         bool `mapstructure:"enableGraphql" json:"enableGraphql"`
}

func (GithubConnection) TableName() string {
	return "_tool_github_connections"
}

// Using GithubUserOfToken because it requires authentication, and it is public information anyway.
type GithubUserOfToken struct {
	Login string `json:"login"`
}
