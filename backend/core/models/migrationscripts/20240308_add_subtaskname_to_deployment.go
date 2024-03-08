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

var _ plugin.MigrationScript = (*addSubtabknameToDeployment)(nil)

type cicdDeployment20240308 struct {
	SubtaskName string `gorm:"type:varchar(255)"`
}

func (cicdDeployment20240308) TableName() string {
	return "cicd_deployments"
}

type cicdDeploymentCommits20240308 struct {
	SubtaskName string `gorm:"type:varchar(255)"`
}

func (cicdDeploymentCommits20240308) TableName() string {
	return "cicd_deployment_commits"
}

type addSubtabknameToDeployment struct{}

func (*addSubtabknameToDeployment) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if err := db.AutoMigrate(&cicdDeployment20240308{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&cicdDeploymentCommits20240308{}); err != nil {
		return err
	}
	return nil
}

func (*addSubtabknameToDeployment) Version() uint64 {
	return 20240308142101
}

func (*addSubtabknameToDeployment) Name() string {
	return "add subtaskname to deployment"
}
