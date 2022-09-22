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
	"encoding/base64"
	"strings"

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/core"
	"gorm.io/gorm"
)

type addInitTables struct{}

func (*addInitTables) Up(ctx context.Context, db *gorm.DB) errors.Error {
	err := db.Migrator().DropTable(
		//history table
		&archived.BitbucketRepo{},
		&archived.BitbucketRepoCommit{},
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

	var result *gorm.DB
	var bitbucketConns []archived.BitbucketConnection
	result = db.Find(&bitbucketConns)
	if result.Error != nil {
		return errors.Convert(result.Error)
	}

	for _, v := range bitbucketConns {
		conn := &archived.BitbucketConnection{}
		conn.ID = v.ID
		conn.Name = v.Name
		conn.Endpoint = v.Endpoint
		conn.Proxy = v.Proxy
		conn.RateLimitPerHour = v.RateLimitPerHour

		c := config.GetConfig()
		encKey := c.GetString(core.EncodeKeyEnvStr)
		if encKey == "" {
			return errors.BadInput.New("bitbucket invalid encKey")
		}
		var auth string
		if auth, err = core.Decrypt(encKey, v.BasicAuth.GetEncodedToken()); err != nil {
			return errors.Convert(err)
		}
		var pk []byte
		pk, err = base64.StdEncoding.DecodeString(auth)
		if err != nil {
			return errors.Default.Wrap(err, "error creating connection entry for BitBucket")
		}
		originInfo := strings.Split(string(pk), ":")
		if len(originInfo) == 2 {
			conn.Username = originInfo[0]
			conn.Password, err = core.Encrypt(encKey, originInfo[1])
			if err != nil {
				return errors.Convert(err)
			}
			// create
			tx := db.Create(&conn)
			if tx.Error != nil {
				return errors.Default.Wrap(tx.Error, "error adding connection to DB")
			}
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
