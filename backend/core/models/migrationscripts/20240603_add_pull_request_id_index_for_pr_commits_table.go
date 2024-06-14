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
	"net/url"
	"strings"
)

type addPullRequestIdIndexToPullRequestCommits struct{}

func (*addPullRequestIdIndexToPullRequestCommits) Up(basicRes context.BasicRes) errors.Error {
	dbUrl := basicRes.GetConfig("DB_URL")
	if dbUrl == "" {
		return errors.BadInput.New("DB_URL is required")
	}
	u, err1 := url.Parse(dbUrl)
	if err1 != nil {
		return errors.Convert(err1)
	}
	switch strings.ToLower(u.Scheme) {
	case "mysql":
		db := basicRes.GetDal()
		err := db.Exec("ALTER TABLE pull_request_commits DROP PRIMARY KEY;")
		if err != nil {
			return err
		}
		err = db.Exec("ALTER TABLE pull_request_commits ADD PRIMARY KEY (pull_request_id, commit_sha);")
		if err != nil {
			return err
		}
		return nil
	case "postgresql", "postgres", "pg":
		db := basicRes.GetDal()
		err := db.Exec("ALTER TABLE pull_request_commits DROP CONSTRAINT pull_request_commits_pkey;")
		if err != nil {
			return err
		}
		err = db.Exec("ALTER TABLE pull_request_commits ADD PRIMARY KEY (pull_request_id, commit_sha);")
		if err != nil {
			return err
		}
		return nil
	default:
		return nil
	}
}

func (*addPullRequestIdIndexToPullRequestCommits) Version() uint64 {
	return 20240602103400
}

func (*addPullRequestIdIndexToPullRequestCommits) Name() string {
	return "changing pull_request_commits primary key columns order"
}
