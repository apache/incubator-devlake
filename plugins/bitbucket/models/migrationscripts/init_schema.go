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
	"github.com/apache/incubator-devlake/plugins/bitbucket/models/migrationscripts/archived"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InitSchemas struct{}

func (*InitSchemas) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().DropTable(
		&archived.BitbucketUser{},
		&archived.BitbucketRepo{},
		&archived.BitbucketRepoCommit{},
		&archived.BitbucketConnection{},
		&archived.BitbucketAccount{},
		&archived.BitbucketCommit{},
		&archived.BitbucketPullRequest{},
		&archived.BitbucketIssue{},
		"_raw_bitbucket_api_repositories",
	)

	if err != nil {
		return err
	}

	err = db.Migrator().AutoMigrate(
		&archived.BitbucketUser{},
		&archived.BitbucketRepo{},
		&archived.BitbucketRepoCommit{},
		&archived.BitbucketConnection{},
		&archived.BitbucketAccount{},
		&archived.BitbucketCommit{},
		&archived.BitbucketPullRequest{},
		&archived.BitbucketIssue{},
	)

	if err != nil {
		return err
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
		if err != nil {
			return err
		}
		conn.Proxy = v.GetString("BITBUCKET_PROXY")
		conn.RateLimitPerHour = v.GetInt("BITBUCKET_API_REQUESTS_PER_HOUR")

		err = db.Clauses(clause.OnConflict{DoNothing: true}).Create(conn).Error

		if err != nil {
			return err
		}
	}

	return nil
}

func (*InitSchemas) Version() uint64 {
	return 20220622200403
}

func (*InitSchemas) Name() string {
	return "Bitbucket init schemas"
}
