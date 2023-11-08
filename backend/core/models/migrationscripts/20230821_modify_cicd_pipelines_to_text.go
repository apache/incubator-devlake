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

var _ plugin.MigrationScript = (*modifyCicdPipelinesToText)(nil)

type modifyCicdPipelinesToText struct{}

type cicdPipeline20230821 struct {
	Name string
}

func (cicdPipeline20230821) TableName() string {
	return "cicd_pipelines"
}

func (script *modifyCicdPipelinesToText) Up(basicRes context.BasicRes) errors.Error {
	// cicd_pipelines.name might be text, we ought to change the type
	// for the column from `varchar(255)` to `text`
	db := basicRes.GetDal()
	return migrationhelper.ChangeColumnsType[cicdPipeline20230821](
		basicRes,
		script,
		cicdPipeline20230821{}.TableName(),
		[]string{"name"},
		func(tmpColumnParams []interface{}) errors.Error {
			return db.UpdateColumn(
				&cicdPipeline20230821{},
				"name",
				dal.DalClause{Expr: " ? ", Params: tmpColumnParams},
				dal.Where("? != '' ", tmpColumnParams...),
			)
		},
	)
}

func (*modifyCicdPipelinesToText) Version() uint64 {
	return 20230821145125
}

func (*modifyCicdPipelinesToText) Name() string {
	return "modify cicd_pipelines name from varchar to text"
}
