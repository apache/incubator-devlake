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

package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*addApiTokenAuth)(nil)

type bitbucketConnection20251001 struct {
	UsesApiToken bool `mapstructure:"usesApiToken" json:"usesApiToken"`
}

func (bitbucketConnection20251001) TableName() string {
	return "_tool_bitbucket_connections"
}

type addApiTokenAuth struct{}

func (script *addApiTokenAuth) Up(basicRes context.BasicRes) errors.Error {
	// Add usesApiToken field to support API token tracking
	// Existing connections will default to false (app password method)
	err := migrationhelper.AutoMigrateTables(basicRes, &bitbucketConnection20251001{})
	if err != nil {
		return err
	}

	// Set default usesApiToken to false for existing connections
	// This ensures backward compatibility with existing App password connections
	db := basicRes.GetDal()
	err = db.Exec("UPDATE _tool_bitbucket_connections SET uses_api_token = false WHERE uses_api_token IS NULL")
	if err != nil {
		return err
	}

	return nil
}

func (*addApiTokenAuth) Version() uint64 {
	return 20251001000001
}

func (script *addApiTokenAuth) Name() string {
	return "add API token authentication support to Bitbucket connections"
}
