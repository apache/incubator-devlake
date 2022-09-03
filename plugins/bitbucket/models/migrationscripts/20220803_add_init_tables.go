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
	"context"
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models/migrationscripts/archived"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type addInitTables struct{}

func (*addInitTables) Up(ctx context.Context, db *gorm.DB) errors.Error {
	err := db.Migrator().DropTable(
		&archived.BitbucketRepo{},
		&archived.BitbucketRepoCommit{},
		&archived.BitbucketConnection{},
		&archived.BitbucketAccount{},
		&archived.BitbucketCommit{},
		&archived.BitbucketPullRequest{},
		&archived.BitbucketIssue{},
		&archived.BitbucketPrComment{},
		&archived.BitbucketIssueComment{},
	)

	if err != nil {
		return errors.Convert(err)
	}

	err = db.Migrator().AutoMigrate(
		&archived.BitbucketRepo{},
		&archived.BitbucketRepoCommit{},
		&archived.BitbucketConnection{},
		&archived.BitbucketAccount{},
		&archived.BitbucketCommit{},
		&archived.BitbucketPullRequest{},
		&archived.BitbucketIssue{},
		&archived.BitbucketPrComment{},
		&archived.BitbucketIssueComment{},
	)

	if err != nil {
		return errors.Convert(err)
	}

	v := config.GetConfig()
	encKey := v.GetString("ENCODE_KEY")
	endPoint := v.GetString("BITBUCKET_ENDPOINT")
	bitbucketUsername := v.GetString("BITBUCKET_AUTH_USERNAME")
	bitbucketAppPassword := v.GetString("BITBUCKET_AUTH_PASSWORD")
	if encKey == "" || endPoint == "" || bitbucketUsername == "" || bitbucketAppPassword == "" {
		return nil
	} else {
		conn := &archived.BitbucketConnection{}
		conn.Name = "init bitbucket connection"
		conn.ID = 1
		conn.Endpoint = endPoint
		conn.BasicAuth.Username = bitbucketUsername
		conn.BasicAuth.Password = bitbucketAppPassword
		conn.Proxy = v.GetString("BITBUCKET_PROXY")
		conn.RateLimitPerHour = v.GetInt("BITBUCKET_API_REQUESTS_PER_HOUR")

		err = db.Clauses(clause.OnConflict{DoNothing: true}).Create(conn).Error

		if err != nil {
			return errors.Default.Wrap(err, "error creating connection entry for BitBucket", errors.AsUserMessage())
		}
	}

	return nil
}

func (*addInitTables) Version() uint64 {
	return 20220803220824
}

func (*addInitTables) Name() string {
	return "Bitbucket init schema 20220803"
}
