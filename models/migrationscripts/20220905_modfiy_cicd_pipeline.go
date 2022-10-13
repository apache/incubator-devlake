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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/core"

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
)

var _ core.MigrationScript = (*modifyCicdPipeline)(nil)

type modifyCicdPipeline struct{}

type CICDPipelineRelationship20220905 struct {
	ParentPipelineId string `gorm:"primaryKey;type:varchar(255)"`
	ChildPipelineId  string `gorm:"primaryKey;type:varchar(255)"`
	archived.NoPKModel
}

func (CICDPipelineRelationship20220905) TableName() string {
	return "cicd_pipeline_relationships"
}

func (*modifyCicdPipeline) Up(basicRes core.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := migrationhelper.DropColumns(basicRes, "cicd_pipelines", "commit_sha", "branch", "repo")
	if err != nil {
		return err
	}
	err = db.RenameColumn("cicd_pipeline_repos", "repo_url", "repo")
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&CICDPipelineRelationship20220905{})
	if err != nil {
		return errors.Convert(err)
	}
	return nil
}

func (*modifyCicdPipeline) Version() uint64 {
	return 20220905232735
}

func (*modifyCicdPipeline) Name() string {
	return "modify cicd pipeline"
}
