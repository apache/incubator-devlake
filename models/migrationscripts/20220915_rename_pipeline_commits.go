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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"gorm.io/gorm"
)

type renamePipelineCommits struct{}

type CiCDPipelineRepoOld struct {
	domainlayer.DomainEntity
	CommitSha string `gorm:"primaryKey;type:varchar(255)"`
	Branch    string `gorm:"type:varchar(255)"`
	Repo      string `gorm:"type:varchar(255)"`
}

func (CiCDPipelineRepoOld) TableName() string {
	return "cicd_pipeline_repos"
}

type CiCDPipelineRepo0915 struct {
	PipelineId string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha  string `gorm:"primaryKey;type:varchar(255)"`
	Branch     string `gorm:"type:varchar(255)"`
	Repo       string `gorm:"type:varchar(255)"`
}

func (CiCDPipelineRepo0915) TableName() string {
	return "cicd_pipeline_commits"
}

func (*renamePipelineCommits) Up(ctx context.Context, db *gorm.DB) errors.Error {
	err := db.Migrator().RenameTable(CiCDPipelineRepoOld{}, CiCDPipelineRepo0915{})
	if err != nil {
		return errors.Convert(err)
	}
	err = db.Migrator().RenameColumn(CiCDPipelineRepo0915{}, `id`, `pipeline_id`)
	if err != nil {
		return errors.Convert(err)
	}
	return nil
}

func (*renamePipelineCommits) Version() uint64 {
	return 20220915000025
}

func (*renamePipelineCommits) Name() string {
	return "UpdateSchemas for renamePipelineCommits"
}
