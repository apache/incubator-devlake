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
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*addOriginalStatusAndResultToDevOpsTables)(nil)

type cicdDeployment2023116 struct {
	OriginalStatus string `gorm:"type:varchar(100)"`
	OriginalResult string `gorm:"type:varchar(100)"`
}

func (cicdDeployment2023116) TableName() string {
	return "cicd_deployments"
}

type cicdDeploymentCommit2023116 struct {
	OriginalStatus string `gorm:"type:varchar(100)"`
	OriginalResult string `gorm:"type:varchar(100)"`
}

func (cicdDeploymentCommit2023116) TableName() string {
	return "cicd_deployment_commits"
}

type cicdPipeline2023116 struct {
	OriginalStatus string `gorm:"type:varchar(100)"`
	OriginalResult string `gorm:"type:varchar(100)"`
}

func (cicdPipeline2023116) TableName() string {
	return "cicd_pipelines"
}

type cicdTask2023116 struct {
	OriginalStatus string `gorm:"type:varchar(100)"`
	OriginalResult string `gorm:"type:varchar(100)"`
}

func (cicdTask2023116) TableName() string {
	return "cicd_tasks"
}

type addOriginalStatusAndResultToDevOpsTables struct{}

func (u *addOriginalStatusAndResultToDevOpsTables) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if err := db.AutoMigrate(&cicdTask2023116{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&cicdPipeline2023116{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&cicdDeploymentCommit2023116{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&cicdDeployment2023116{}); err != nil {
		return err
	}
	return nil
}

func (*addOriginalStatusAndResultToDevOpsTables) Version() uint64 {
	return 20231116142100
}

func (*addOriginalStatusAndResultToDevOpsTables) Name() string {
	return "add original status and original result to all related devops tables"
}
