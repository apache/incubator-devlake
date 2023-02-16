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

var _ plugin.MigrationScript = (*modifyIssueStorypointToFloat64)(nil)

type modifyIssueStorypointToFloat64 struct{}

type issues20230214 struct {
	StoryPoint float64
}

func (issues20230214) TableName() string {
	return "issues"
}

func (script *modifyIssueStorypointToFloat64) Up(basicRes context.BasicRes) errors.Error {
	// issues.story_point might be float, we ought to change the type
	// for the column from `int64` to `float64`
	db := basicRes.GetDal()
	return migrationhelper.ChangeColumnsType[issues20230214](
		basicRes,
		script,
		issues20230214{}.TableName(),
		[]string{"story_point"},
		func(tmpColumnParams []interface{}) errors.Error {
			return db.UpdateColumn(
				&issues20230214{},
				"story_point",
				dal.DalClause{Expr: " ? ", Params: tmpColumnParams},
				dal.Where("? is not null ", tmpColumnParams...),
			)
		},
	)
}

func (*modifyIssueStorypointToFloat64) Version() uint64 {
	return 20230214145125
}

func (*modifyIssueStorypointToFloat64) Name() string {
	return "modify issues story_point from int64 to float64"
}
