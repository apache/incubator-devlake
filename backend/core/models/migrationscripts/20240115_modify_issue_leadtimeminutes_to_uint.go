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

var _ plugin.MigrationScript = (*modifyIssueLeadTimeMinutesToUint)(nil)

type issue20240115 struct {
	LeadTimeMinutes *uint
}

func (issue20240115) TableName() string {
	return "issues"
}

type modifyIssueLeadTimeMinutesToUint struct{}

func (u *modifyIssueLeadTimeMinutesToUint) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if err := migrationhelper.ChangeColumnsType[issue20240115](
		basicRes,
		u,
		issue20240115{}.TableName(),
		[]string{"lead_time_minutes"},
		func(tmpColumnParams []interface{}) errors.Error {
			return db.UpdateColumn(
				&issue20240115{},
				"lead_time_minutes",
				dal.DalClause{Expr: " ? ", Params: tmpColumnParams},
				dal.Where("? != 0", tmpColumnParams...),
			)
		},
	); err != nil {
		return err
	}

	return nil
}

func (*modifyIssueLeadTimeMinutesToUint) Version() uint64 {
	return 20240115170000
}

func (*modifyIssueLeadTimeMinutesToUint) Name() string {
	return "modify issues lead_time_minutes to *uint"
}
