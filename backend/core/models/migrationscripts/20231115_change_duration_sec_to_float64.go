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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*changeDurationSecToFloat64)(nil)

type cicdDeployment2023115 struct {
	DurationSec float64
}

func (cicdDeployment2023115) TableName() string {
	return "cicd_deployments"
}

type cicdDeploymentCommit2023115 struct {
	DurationSec float64
}

func (cicdDeploymentCommit2023115) TableName() string {
	return "cicd_deployment_commits"
}

type cicdPipeline2023115 struct {
	DurationSec float64
}

func (cicdPipeline2023115) TableName() string {
	return "cicd_pipelines"
}

type cicdTask2023115 struct {
	DurationSec float64
}

func (cicdTask2023115) TableName() string {
	return "cicd_tasks"
}

type changeDurationSecToFloat64 struct{}

func (u *changeDurationSecToFloat64) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if err := migrationhelper.ChangeColumnsType[cicdDeployment2023115](
		basicRes,
		u,
		cicdDeployment2023115{}.TableName(),
		[]string{"duration_sec"},
		func(tmpColumnParams []interface{}) errors.Error {
			return db.UpdateColumn(
				&cicdDeployment2023115{},
				"duration_sec",
				dal.DalClause{Expr: " ? ", Params: tmpColumnParams},
				dal.Where("? != 0", tmpColumnParams...),
			)
		},
	); err != nil {
		return err
	}
	if err := migrationhelper.ChangeColumnsType[cicdDeploymentCommit2023115](
		basicRes,
		u,
		cicdDeploymentCommit2023115{}.TableName(),
		[]string{"duration_sec"},
		func(tmpColumnParams []interface{}) errors.Error {
			return db.UpdateColumn(
				&cicdDeploymentCommit2023115{},
				"duration_sec",
				dal.DalClause{Expr: " ? ", Params: tmpColumnParams},
				dal.Where("? != 0", tmpColumnParams...),
			)
		},
	); err != nil {
		return err
	}
	if err := migrationhelper.ChangeColumnsType[cicdPipeline2023115](
		basicRes,
		u,
		cicdPipeline2023115{}.TableName(),
		[]string{"duration_sec"},
		func(tmpColumnParams []interface{}) errors.Error {
			return db.UpdateColumn(
				&cicdPipeline2023115{},
				"duration_sec",
				dal.DalClause{Expr: " ? ", Params: tmpColumnParams},
				dal.Where("? != 0", tmpColumnParams...),
			)
		},
	); err != nil {
		return err
	}
	if err := migrationhelper.ChangeColumnsType[cicdTask2023115](
		basicRes,
		u,
		cicdTask2023115{}.TableName(),
		[]string{"duration_sec"},
		func(tmpColumnParams []interface{}) errors.Error {
			return db.UpdateColumn(
				&cicdTask2023115{},
				"duration_sec",
				dal.DalClause{Expr: " ? ", Params: tmpColumnParams},
				dal.Where("? != 0", tmpColumnParams...),
			)
		},
	); err != nil {
		return err
	}
	return nil
}

func (*changeDurationSecToFloat64) Version() uint64 {
	return 20231115170000
}

func (*changeDurationSecToFloat64) Name() string {
	return "change duration_sec field to float64 in all related tables"
}
