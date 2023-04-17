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
	"github.com/apache/incubator-devlake/plugins/bitbucket/models/migrationscripts/archived"
)

type addRepoIdField20230411 struct{}

func (*addRepoIdField20230411) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(&archived.BitbucketPipelineStep20230411{})
}

func (*addRepoIdField20230411) Version() uint64 {
	return 2023041111938
}

func (*addRepoIdField20230411) Name() string {
	return "add column `repo_id` and `commit_sha` at _tool_bitbucket_pipelines"
}
