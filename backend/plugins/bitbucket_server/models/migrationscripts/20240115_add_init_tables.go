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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models/migrationscripts/archived"
)

type addInitTables20240115 struct{}

func (script *addInitTables20240115) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&archived.BitbucketServerUser{},
		&archived.BitbucketServerConnection{},
		&archived.BitbucketServerPullRequest{},
		&archived.BitbucketServerPrComment{},
		&archived.BitbucketServerPrCommit{},
		&archived.BitbucketServerRepo{},
		&archived.BitbucketServerScopeConfig{},
	)
}

func (*addInitTables20240115) Version() uint64 {
	return 20240115
}

func (*addInitTables20240115) Name() string {
	return "Bitbucket Server init schema 20240115"
}
