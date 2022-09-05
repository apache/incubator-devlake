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
	"gorm.io/gorm"
)

type modifyPipeline struct{}

func (*modifyPipeline) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().DropColumn(CICDPipeline0905{}, "commit_sha")
	if err != nil {
		return err
	}
	err = db.Migrator().DropColumn(CICDPipeline0905{}, "branch")
	if err != nil {
		return err
	}
	err = db.Migrator().DropColumn(CICDPipeline0905{}, "repo")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(CICDPipelineRepo0905{}, "repo_url", "repo")
	if err != nil {
		return err
	}
	return nil
}

func (*modifyPipeline) Version() uint64 {
	return 20220905232735
}

func (*modifyPipeline) Name() string {
	return "modify cicd pipeline"
}

type CICDPipeline0905 struct {
	CommitSha string `gorm:"type:varchar(255);index"`
	Branch    string `gorm:"type:varchar(255);index"`
	Repo      string `gorm:"type:varchar(255);index"`
}

func (CICDPipeline0905) TableName() string {
	return "cicd_pipelines"
}

type CICDPipelineRepo0905 struct {
	RepoUrl string `gorm:"type:varchar(255)"`
}

func (CICDPipelineRepo0905) TableName() string {
	return "cicd_pipeline_repos"
}
