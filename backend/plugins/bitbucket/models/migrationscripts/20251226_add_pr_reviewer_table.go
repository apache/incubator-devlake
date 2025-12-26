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

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*addPrReviewerTable)(nil)

type prReviewer20251226 struct {
	ConnectionId   uint64 `gorm:"primaryKey"`
	RepoId         string `gorm:"primaryKey;type:varchar(255)"`
	PullRequestId  int    `gorm:"primaryKey"`
	AccountId      string `gorm:"primaryKey;type:varchar(255)"`
	DisplayName    string `gorm:"type:varchar(255)"`
	Role           string `gorm:"type:varchar(100)"`
	State          string `gorm:"type:varchar(100)"`
	Approved       bool
	ParticipatedOn *time.Time
	archived.NoPKModel
}

func (prReviewer20251226) TableName() string {
	return "_tool_bitbucket_pr_reviewers"
}

type addPrReviewerTable struct{}

func (*addPrReviewerTable) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if err := db.AutoMigrate(&prReviewer20251226{}); err != nil {
		return err
	}
	return nil
}

func (*addPrReviewerTable) Version() uint64 {
	return 20251226100000
}

func (*addPrReviewerTable) Name() string {
	return "add _tool_bitbucket_pr_reviewers table"
}
