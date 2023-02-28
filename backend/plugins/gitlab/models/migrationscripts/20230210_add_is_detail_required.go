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
	"time"
)

type gitlabMergeRequests20230210 struct {
	IsDetailRequired bool
	GitlabUpdatedAt  *time.Time
}

func (gitlabMergeRequests20230210) TableName() string {
	return "_tool_gitlab_merge_requests"
}

type gitlabPipelines20230210 struct {
	IsDetailRequired bool
}

func (gitlabPipelines20230210) TableName() string {
	return "_tool_gitlab_pipelines"
}

type addIsDetailRequired20230210 struct{}

func (script *addIsDetailRequired20230210) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&gitlabMergeRequests20230210{},
		&gitlabPipelines20230210{},
	)
}

func (*addIsDetailRequired20230210) Version() uint64 {
	return 20230228161031
}

func (*addIsDetailRequired20230210) Name() string {
	return "add IsDetailRequired to _tool_gitlab_merge_requests and _tool_gitlab_pipelines"
}
