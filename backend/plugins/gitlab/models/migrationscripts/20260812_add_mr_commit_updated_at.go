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
	"time"

	"github.com/apache/devlake/core/context"
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
)

var _ plugin.MigrationScript = (*addMrCommitUpdatedAt)(nil)

type mrCommitUpdatedAt struct {
	CommitUpdatedAt *time.Time
}

func (mrCommitUpdatedAt) TableName() string {
	return "_tool_gitlab_merge_requests"
}

type addMrCommitUpdatedAt struct{}

func (*addMrCommitUpdatedAt) Up(basicRes context.BasicRes) errors.Error {
	return errors.Convert(basicRes.GetDal().AutoMigrate(&mrCommitUpdatedAt{}))
}

func (*addMrCommitUpdatedAt) Version() uint64 {
	return 20260812000001
}

func (*addMrCommitUpdatedAt) Name() string {
	return "gitlab: add commit_updated_at to merge requests"
}
