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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/plugins/core"
)

var _ core.MigrationScript = (*commitLineChange)(nil)

type commitLineChange struct{}

type commitLineChange20220918 struct {
	domainlayer.DomainEntity
	Id          string `gorm:"type:varchar(255);primaryKey"`
	CommitSha   string `gorm:"type:varchar(40);"`
	NewFilePath string `gorm:"type:varchar(255);"`
	LineNoNew   int    `gorm:"type:int"`
	LineNoOld   int    `gorm:"type:int"`
	OldFilePath string `gorm:"type:varchar(255)"`
	HunkNum     int    `gorm:"type:int"`
	ChangedType string `gorm:"type:varchar(255)"`
	PrevCommit  string `gorm:"type:varchar(255)"`
}

func (commitLineChange20220918) TableName() string {
	return "commit_line_change"
}

func (*commitLineChange) Up(basicRes core.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(&commitLineChange20220918{})

}

func (*commitLineChange) Version() uint64 {
	return 202209221033
}

func (*commitLineChange) Name() string {

	return "add commit line change table"
}
