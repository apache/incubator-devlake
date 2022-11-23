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
	"strconv"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/models/migrationscripts/archived"
)

type addInitTables struct{}

func (*addInitTables) Up(baseRes core.BasicRes) errors.Error {
	db := baseRes.GetDal()
	err := db.DropTables(
		&archived.GitlabProject{},
		&archived.GitlabMergeRequest{},
		&archived.GitlabCommit{},
		&archived.GitlabTag{},
		&archived.GitlabProjectCommit{},
		&archived.GitlabPipeline{},
		&archived.GitlabReviewer{},
		&archived.GitlabMrNote{},
		&archived.GitlabMrCommit{},
		&archived.GitlabMrComment{},
		&archived.GitlabConnection{},
		&archived.GitlabIssue{},
		&archived.GitlabIssueLabel{},
		&archived.GitlabMrLabel{},
		"_tool_gitlab_users",
		"_raw_gitlab_api_children_on_pipeline",
		"_raw_gitlab_api_commit",
		"_raw_gitlab_api_issues",
		"_raw_gitlab_api_merge_request_commits",
		"_raw_gitlab_api_merge_request_notes",
		"_raw_gitlab_api_merge_requests",
		"_tool_gitlab_merge_request_comments",
		"_tool_gitlab_merge_request_commits",
		"_tool_gitlab_merge_request_notes",
		"_raw_gitlab_api_pipeline",
		"_raw_gitlab_api_project",
		"_raw_gitlab_api_tag",
	)

	if err != nil {
		return err
	}

	err = migrationhelper.AutoMigrateTables(
		baseRes,
		&archived.GitlabProject{},
		&archived.GitlabMergeRequest{},
		&archived.GitlabCommit{},
		&archived.GitlabTag{},
		&archived.GitlabProjectCommit{},
		&archived.GitlabPipeline{},
		&archived.GitlabReviewer{},
		&archived.GitlabMrNote{},
		&archived.GitlabMrCommit{},
		&archived.GitlabMrComment{},
		&archived.GitlabAccount{},
		&archived.GitlabConnection{},
		&archived.GitlabIssue{},
		&archived.GitlabIssueLabel{},
		&archived.GitlabMrLabel{},
	)

	if err != nil {
		return err
	}

	encKey := baseRes.GetConfig("ENCODE_KEY")
	endPoint := baseRes.GetConfig("GITLAB_ENDPOINT")
	gitlabAuth := baseRes.GetConfig("GITLAB_AUTH")

	if encKey == "" || endPoint == "" || gitlabAuth == "" {
		return nil
	}
	conn := &archived.GitlabConnection{}
	conn.Name = "init gitlab connection"
	conn.ID = 1
	conn.Endpoint = endPoint
	conn.Token, err = core.Encrypt(encKey, gitlabAuth)
	if err != nil {
		return err
	}
	conn.Proxy = baseRes.GetConfig("GITLAB_PROXY")
	var err1 error
	conn.RateLimitPerHour, err1 = strconv.Atoi(baseRes.GetConfig("GITLAB_API_REQUESTS_PER_HOUR"))
	if err1 != nil {
		conn.RateLimitPerHour = 1000
	}
	err = db.CreateIfNotExist(conn)

	if err != nil {
		return err
	}

	return nil
}

func (*addInitTables) Version() uint64 {
	return 20220714231236
}

func (*addInitTables) Name() string {
	return "Gitlab init schemas"
}
