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
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*modifyRefsIdLength)(nil)

type modifyRefsIdLength struct{}
type ref20240124 struct {
	archived.DomainEntityExtended
	RepoId      string `gorm:"type:varchar(255)"`
	Name        string `gorm:"type:varchar(255)"`
	CommitSha   string `gorm:"type:varchar(40)"`
	IsDefault   bool
	RefType     string `gorm:"type:varchar(255)"`
	CreatedDate *time.Time
}

func (ref20240124) TableName() string {
	return "refs"
}

func (script *modifyRefsIdLength) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	return migrationhelper.ChangePrimaryKeyColumnsType[ref20240124](
		basicRes,
		script,
		ref20240124{}.TableName(),
		[]string{"id"},
		func(tmpColumnParams []interface{}) errors.Error {
			return db.UpdateColumn(
				&ref20240124{},
				"id",
				dal.DalClause{Expr: " ? ", Params: tmpColumnParams},
				dal.Where("? IS NOT NULL", tmpColumnParams...),
			)
		},
	)
}

func (*modifyRefsIdLength) Version() uint64 {
	return 20240124155129
}

func (*modifyRefsIdLength) Name() string {
	return "modify refs id length from 255 to 500"
}
