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
	"github.com/apache/incubator-devlake/plugins/gitlab/models/migrationscripts/archived"
)

type addGitlabAssignee struct{}

func (*addGitlabAssignee) Up(baseRes context.BasicRes) errors.Error {
	err := baseRes.GetDal().DropTables(archived.GitlabAssignee{})
	if err != nil {
		return err
	}

	err = migrationhelper.AutoMigrateTables(
		baseRes,
		&archived.GitlabAssignee{},
		&archived.GitlabReviewer{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (*addGitlabAssignee) Version() uint64 {
	return 20240531110339
}

func (*addGitlabAssignee) Name() string {
	return "add _tool_gitlab_assignees table"
}
