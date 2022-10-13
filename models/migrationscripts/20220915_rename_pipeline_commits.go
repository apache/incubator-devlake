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
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/core"
)

var _ core.MigrationScript = (*renamePipelineCommits)(nil)

type renamePipelineCommits struct{}

type CiCDPipelineRepo20220915Before struct {
	archived.DomainEntity
	CommitSha string `gorm:"primaryKey;type:varchar(255)"`
	Branch    string `gorm:"type:varchar(255)"`
	Repo      string `gorm:"type:varchar(255)"`
}

func (CiCDPipelineRepo20220915Before) TableName() string {
	return "cicd_pipeline_repos"
}

type CiCDPipelineRepo20220915After struct {
	archived.NoPKModel
	PipelineId string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha  string `gorm:"primaryKey;type:varchar(255)"`
	Branch     string `gorm:"type:varchar(255)"`
	RepoId     string `gorm:"index;type:varchar(255)"`
	RepoUrl    string
}

func (CiCDPipelineRepo20220915After) TableName() string {
	return "cicd_pipeline_commits"
}

func (*renamePipelineCommits) Up(basicRes core.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.RenameTable("cicd_pipeline_repos", "cicd_pipeline_commits")
	if err != nil {
		return err
	}
	// err = db.DropIndexes("cicd_pipeline_repos", `idx_cicd_pipeline_repos_raw_data_params`)
	// if err != nil {
	// 	return err
	// }
	err = db.RenameColumn("cicd_pipeline_commits", "id", "pipeline_id")
	if err != nil {
		return err
	}
	err = db.AutoMigrate(CiCDPipelineRepo20220915After{})
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
