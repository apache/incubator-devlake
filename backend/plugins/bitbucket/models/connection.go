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
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

var _ plugin.ApiConnection = (*BitbucketConnection)(nil)

// BitbucketConn holds the essential information to connect to the Bitbucket API
type BitbucketConn struct {
	api.RestConnection `mapstructure:",squash"`
	api.BasicAuth      `mapstructure:",squash"`
	// UsesApiToken indicates whether the password field contains an API token (true)
	// or an App password (false). Both use Basic Auth, but API tokens are the new standard.
	UsesApiToken bool `mapstructure:"usesApiToken" json:"usesApiToken"`
}

func (bc BitbucketConn) Sanitize() BitbucketConn {
	bc.Password = ""
	return bc
}

// SetupAuthentication sets up HTTP Basic Authentication
// Both App passwords and API tokens use Basic Auth with username:credential format
func (bc *BitbucketConn) SetupAuthentication(req *http.Request) errors.Error {
	return bc.BasicAuth.SetupAuthentication(req)
}

// BitbucketConnection holds BitbucketConn plus ID/Name for database storage
type BitbucketConnection struct {
	api.BaseConnection `mapstructure:",squash"`
	BitbucketConn      `mapstructure:",squash"`
}

func (BitbucketConnection) TableName() string {
	return "_tool_bitbucket_connections"
}

func (connection BitbucketConnection) Sanitize() BitbucketConnection {
	connection.BitbucketConn = connection.BitbucketConn.Sanitize()
	return connection
}

func (connection *BitbucketConnection) MergeFromRequest(target *BitbucketConnection, body map[string]interface{}) error {
	password := target.Password
	if err := api.DecodeMapStruct(body, target, true); err != nil {
		return err
	}
	modifiedPassword := target.Password
	if modifiedPassword == "" {
		target.Password = password
	}
	return nil
}
