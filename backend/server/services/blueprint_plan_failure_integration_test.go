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

package services

import (
	"os"
	"testing"
	"time"

	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/models"
	contextimpl "github.com/apache/devlake/impls/context"
	"github.com/apache/devlake/impls/dalgorm"
	"github.com/apache/devlake/impls/logruslog"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Use dedicated, disposable databases: this test creates the pipeline tables.
// Unset DSNs keep ordinary unit-test runs independent of database servers.
func TestBlueprintPlanFailureDatabase(t *testing.T) {
	for _, driver := range []string{"mysql", "postgres"} {
		t.Run(driver, func(t *testing.T) {
			env := "DEVLAKE_TEST_MYSQL_DSN"
			if driver == "postgres" {
				env = "DEVLAKE_TEST_POSTGRES_DSN"
			}
			dsn := os.Getenv(env)
			if dsn == "" {
				t.Skip("set " + env + " to a disposable database")
			}
			var dialect gorm.Dialector = mysql.Open(dsn)
			if driver == "postgres" {
				dialect = postgres.Open(dsn)
			}
			dalgorm.Init("01234567890123456789012345678901")
			database, err := gorm.Open(dialect, &gorm.Config{})
			require.NoError(t, err)
			sqlDB, err := database.DB()
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, sqlDB.Close()) })
			require.NoError(t, database.AutoMigrate(&models.Pipeline{}, &models.DbPipelineLabel{}, &models.Task{}))
			testDB := dalgorm.NewDalgorm(database)
			withProjectTestDatabase(t, testDB)
			previousRes := basicRes
			basicRes = contextimpl.NewDefaultBasicRes(nil, logruslog.Global, testDB)
			t.Cleanup(func() { basicRes = previousRes })

			bp := &models.Blueprint{Name: "planning-failure-integration", Labels: []string{"integration"}}
			bp.ID = uint64(time.Now().UnixNano())
			t.Cleanup(func() {
				var ids []uint64
				require.NoError(t, database.Model(&models.Pipeline{}).Where("blueprint_id = ?", bp.ID).Pluck("id", &ids).Error)
				if len(ids) > 0 {
					require.NoError(t, database.Where("pipeline_id IN ?", ids).Delete(&models.DbPipelineLabel{}).Error)
					require.NoError(t, database.Where("blueprint_id = ?", bp.ID).Delete(&models.Pipeline{}).Error)
				}
			})
			planErr := errors.Default.New("source planning failed")
			stored, saveErr := createFailedBlueprintPipeline(bp, time.Now(), planErr)
			require.NoError(t, saveErr)
			require.NotZero(t, stored.ID)
			rows, count, readErr := GetDbPipelines(&PipelineQuery{BlueprintId: bp.ID, Status: models.TASK_FAILED, Label: "integration"})
			require.NoError(t, readErr)
			require.EqualValues(t, 1, count)
			require.Len(t, rows, 1)
			require.Equal(t, stored.ID, rows[0].ID)
			require.Equal(t, bp.Labels, rows[0].Labels)
			require.Empty(t, rows[0].Plan)
			require.Contains(t, rows[0].Message, planErr.Error())
			require.NotNil(t, rows[0].BeganAt)
			require.NotNil(t, rows[0].FinishedAt)
			_, pending, readErr := GetDbPipelines(&PipelineQuery{BlueprintId: bp.ID, Pending: 1})
			require.NoError(t, readErr)
			require.Zero(t, pending)
			var tasks int64
			require.NoError(t, database.Model(&models.Task{}).Where("pipeline_id = ?", stored.ID).Count(&tasks).Error)
			require.Zero(t, tasks)

			// The duplicate primary key makes the label INSERT fail after the
			// pipeline INSERT. Neither the second pipeline nor its labels may survive.
			bp.Labels = []string{"duplicate", "duplicate"}
			failed, saveErr := createFailedBlueprintPipeline(bp, time.Now(), planErr)
			require.Error(t, saveErr)
			require.Nil(t, failed)
			var pipelines, labels int64
			require.NoError(t, database.Model(&models.Pipeline{}).Where("blueprint_id = ?", bp.ID).Count(&pipelines).Error)
			require.EqualValues(t, 1, pipelines)
			require.NoError(t, database.Model(&models.DbPipelineLabel{}).Where("name = ?", "duplicate").Count(&labels).Error)
			require.Zero(t, labels)
		})
	}
}
