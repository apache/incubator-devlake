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

var _ plugin.MigrationScript = (*modifyPrLabelsAndComments)(nil)

type modifyPrLabelsAndComments struct{}

type pullRequestComment20230516 struct {
	PullRequestId string `gorm:"index;type:varchar(255)"`
}

func (pullRequestComment20230516) TableName() string {
	return "pull_request_comments"
}

func (script *modifyPrLabelsAndComments) Up(basicRes context.BasicRes) errors.Error {
	// change pull_request_comments.pr_id to varchar(255) type
	db := basicRes.GetDal()
	return migrationhelper.ChangeColumnsType[pullRequestComment20230516](
		basicRes,
		script,
		pullRequestComment20230516{}.TableName(),
		[]string{"pull_request_id"},
		func(tmpColumnParams []interface{}) errors.Error {
			return db.UpdateColumn(
				&pullRequestComment20230516{},
				"pull_request_id",
				dal.DalClause{Expr: " ? ", Params: tmpColumnParams},
				dal.Where("? is not null", tmpColumnParams...),
			)
		},
	)
}

func (*modifyPrLabelsAndComments) Version() uint64 {
	return 20230516000001
}

func (*modifyPrLabelsAndComments) Name() string {
	return "change pull_request_comments.pr_id to varchar(255) type"
}
