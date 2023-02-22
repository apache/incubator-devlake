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
)

var _ plugin.MigrationScript = (*addCommitShaIndex)(nil)

type addCommitShaIndex struct{}

func (script *addCommitShaIndex) Up(basicRes context.BasicRes) errors.Error {

	db := basicRes.GetDal()
	err := db.Exec("CREATE INDEX idx_commit_files_commit_sha ON commit_files (commit_sha)")
	if err != nil {
		return err
	}
	err = db.Exec("CREATE INDEX idx_repo_commits_commit_sha ON repo_commits (commit_sha);")
	if err != nil {
		return err
	}
	return nil
}

func (*addCommitShaIndex) Version() uint64 {
	return 20230222145125
}

func (*addCommitShaIndex) Name() string {
	return "add commit_sha index for commit_files and repo_commits"
}
