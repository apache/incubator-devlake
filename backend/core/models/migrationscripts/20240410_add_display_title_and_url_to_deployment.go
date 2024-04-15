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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*addDisplayTitleAndUrl)(nil)

type cicdDeployment20240410 struct {
	DisplayTitle string
	Url          string
}

func (cicdDeployment20240410) TableName() string {
	return "cicd_deployments"
}

type cicdDeploymentCommit20240410 struct {
	DisplayTitle string
	Url          string
}

func (cicdDeploymentCommit20240410) TableName() string {
	return "cicd_deployment_commits"
}

type cicdPipeline20240410 struct {
	DisplayTitle string
	Url          string
}

func (cicdPipeline20240410) TableName() string {
	return "cicd_pipelines"
}

type cicdPipelineCommit20240410 struct {
	DisplayTitle string
	Url          string
}

func (cicdPipelineCommit20240410) TableName() string {
	return "cicd_pipeline_commits"
}

type addDisplayTitleAndUrl struct{}

func (*addDisplayTitleAndUrl) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &cicdDeployment20240410{}, &cicdDeploymentCommit20240410{}, &cicdPipeline20240410{}, &cicdPipelineCommit20240410{})
}

func (*addDisplayTitleAndUrl) Version() uint64 {
	return 20240410111248
}

func (*addDisplayTitleAndUrl) Name() string {
	return "add display title and url to deployments and pipelines"
}
