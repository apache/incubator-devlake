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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

var _ plugin.ApiConnection = (*BitbucketConnection)(nil)

// BitbucketConn holds the essential information to connect to the Bitbucket API
type BitbucketConn struct {
	api.RestConnection `mapstructure:",squash"`
	api.BasicAuth      `mapstructure:",squash"`
}

// BitbucketConnection holds BitbucketConn plus ID/Name for database storage
type BitbucketConnection struct {
	api.BaseConnection `mapstructure:",squash"`
	BitbucketConn      `mapstructure:",squash"`
}

func (BitbucketConnection) TableName() string {
	return "_tool_bitbucket_connections"
}
