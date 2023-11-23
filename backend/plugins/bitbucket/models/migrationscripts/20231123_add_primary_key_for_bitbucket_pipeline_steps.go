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
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
	"time"
)

var _ plugin.MigrationScript = (*reCreatBitBucketPipelineSteps)(nil)

type bitbucketPipelineStep20231123 struct {
	ConnectionId      uint64 `gorm:"primaryKey"`
	BitbucketId       string `gorm:"primaryKey"`
	PipelineId        string `gorm:"type:varchar(255)"`
	Name              string `gorm:"type:varchar(255)"`
	Trigger           string `gorm:"type:varchar(255)"`
	State             string `gorm:"type:varchar(255)"`
	Result            string `gorm:"type:varchar(255)"`
	RepoId            string `gorm:"type:varchar(255)"`
	MaxTime           int
	StartedOn         *time.Time
	CompletedOn       *time.Time
	DurationInSeconds int
	BuildSecondsUsed  int
	RunNumber         int
	Type              string `gorm:"type:varchar(255)"`
	Environment       string `gorm:"type:varchar(255)"`
	archived.NoPKModel
}

func (bitbucketPipelineStep20231123) TableName() string {
	return "_tool_bitbucket_pipeline_steps"
}

type reCreatBitBucketPipelineSteps struct{}

func (script *reCreatBitBucketPipelineSteps) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if err := db.DropTables(bitbucketPipelineStep20231123{}.TableName()); err != nil {
		return err
	}
	return db.AutoMigrate(&bitbucketPipelineStep20231123{})
}

func (*reCreatBitBucketPipelineSteps) Version() uint64 {
	return 20231123160001
}

func (script *reCreatBitBucketPipelineSteps) Name() string {
	return "re create _tool_bitbucket_pipeline_steps, make sure primary keys exist."
}
