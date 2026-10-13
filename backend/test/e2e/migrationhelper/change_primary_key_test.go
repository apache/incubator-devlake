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

package migrationhelper

import (
	"fmt"
	"testing"

	"github.com/apache/devlake/core/config"
	"github.com/apache/devlake/core/context"
	"github.com/apache/devlake/core/dal"
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
	"github.com/apache/devlake/core/runner"
	"github.com/apache/devlake/helpers/e2ehelper"
	"github.com/apache/devlake/helpers/migrationhelper"
	"github.com/apache/devlake/impls/logruslog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const changePkTable = "change_pk_test"

type changePkBefore struct {
	Id   string `gorm:"primaryKey;type:varchar(255)"`
	Name string `gorm:"type:varchar(255)"`
}

func (changePkBefore) TableName() string {
	return changePkTable
}

type changePkAfter struct {
	Id   string `gorm:"primaryKey;type:varchar(500)"`
	Name string `gorm:"type:varchar(255)"`
}

func (changePkAfter) TableName() string {
	return changePkTable
}

type changePkScript struct{}

func (*changePkScript) Up(context.BasicRes) errors.Error { return nil }
func (*changePkScript) Version() uint64                  { return 20261006000000 }
func (*changePkScript) Name() string                     { return "change primary key test" }

var _ plugin.MigrationScript = (*changePkScript)(nil)

// TestChangePrimaryKeyColumnsType checks that the table never loses its primary
// key: MySQL Group Replication rejects DML on tables without one
// (https://github.com/apache/devlake/issues/9194). On MySQL the check is done by
// the database itself through sql_require_primary_key, which rejects every
// ALTER TABLE that would leave the table without a primary key.
//
// Requires E2E_DB_URL (runs under `make e2e-test`).
func TestChangePrimaryKeyColumnsType(t *testing.T) {
	for _, dropPrimaryKeyFirst := range []bool{false, true} {
		name := "with primary key"
		if dropPrimaryKeyFirst {
			// the state left behind by earlier versions of the helper
			name = "without primary key"
		}
		t.Run(name, func(t *testing.T) {
			db := e2ehelper.NewIsolatedMigrationDb(t, "change_pk")
			basicRes := runner.CreateBasicRes(config.GetConfig(), logruslog.Global, db)
			dalInstance := basicRes.GetDal()

			require.NoError(t, dalInstance.AutoMigrate(&changePkBefore{}))
			require.NoError(t, dalInstance.CreateOrUpdate([]*changePkBefore{
				{Id: "a", Name: "first"},
				{Id: "b", Name: "second"},
			}))
			if dropPrimaryKeyFirst {
				dropPrimaryKey := fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s_pkey", changePkTable, changePkTable)
				if db.Dialector.Name() == "mysql" {
					dropPrimaryKey = fmt.Sprintf("ALTER TABLE %s DROP PRIMARY KEY", changePkTable)
				}
				require.NoError(t, dalInstance.Exec(dropPrimaryKey))
			}
			if db.Dialector.Name() == "mysql" {
				// sql_require_primary_key is a session variable, keep a single
				// connection so it applies to every statement of the helper
				sqlDb, err := db.DB()
				require.NoError(t, err)
				sqlDb.SetMaxOpenConns(1)
				require.NoError(t, dalInstance.Exec("SET SESSION sql_require_primary_key = ON"))
			}

			err := migrationhelper.ChangePrimaryKeyColumnsType[changePkAfter](
				basicRes,
				&changePkScript{},
				changePkTable,
				[]string{"id"},
				func(tmpColumnParams []interface{}) errors.Error {
					primaryKeyColumns, err := dal.GetPrimarykeyColumns(dalInstance, dal.DefaultTabler{Name: changePkTable})
					require.NoError(t, err)
					assert.NotEmpty(t, primaryKeyColumns, "table must keep a primary key while data is copied")
					return dalInstance.UpdateColumn(
						&changePkAfter{},
						"id",
						dal.DalClause{Expr: " ? ", Params: tmpColumnParams},
						dal.Where("? IS NOT NULL", tmpColumnParams...),
					)
				},
			)
			require.NoError(t, err)

			primaryKeyColumns, err := dal.GetPrimarykeyColumns(dalInstance, dal.DefaultTabler{Name: changePkTable})
			require.NoError(t, err)
			require.Len(t, primaryKeyColumns, 1)
			assert.Equal(t, "id", primaryKeyColumns[0].Name())
			length, ok := primaryKeyColumns[0].Length()
			assert.True(t, ok)
			assert.Equal(t, int64(500), length)

			columnNames, err := dal.GetColumnNames(dalInstance, dal.DefaultTabler{Name: changePkTable}, func(dal.ColumnMeta) bool { return true })
			require.NoError(t, err)
			assert.ElementsMatch(t, []string{"id", "name"}, columnNames)

			var rows []changePkAfter
			require.NoError(t, dalInstance.All(&rows, dal.Orderby("id")))
			assert.Equal(t, []changePkAfter{{Id: "a", Name: "first"}, {Id: "b", Name: "second"}}, rows)
		})
	}
}
