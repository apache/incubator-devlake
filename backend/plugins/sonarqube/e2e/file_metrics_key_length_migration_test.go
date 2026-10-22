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

package e2e

import (
	"testing"

	"github.com/apache/devlake/core/models/common"
	"github.com/apache/devlake/helpers/e2ehelper"
	implcontext "github.com/apache/devlake/impls/context"
	"github.com/apache/devlake/plugins/sonarqube/impl"
	sonarqubemigrations "github.com/apache/devlake/plugins/sonarqube/models/migrationscripts"
	"github.com/stretchr/testify/require"
)

type sonarqubeFileMetricsBeforeKeyLength struct {
	ConnectionId   uint64 `gorm:"primaryKey"`
	FileMetricsKey string `gorm:"primaryKey;type:varchar(191)"`
	ProjectKey     string `gorm:"index"`
	common.NoPKModel
}

func (sonarqubeFileMetricsBeforeKeyLength) TableName() string {
	return "_tool_sonarqube_file_metrics"
}

// TestSonarqubeFileMetricsKeyLengthMigrationRerun covers the normal run and a
// re-run after an interrupted attempt that already dropped the primary key.
func TestSonarqubeFileMetricsKeyLengthMigrationRerun(t *testing.T) {
	var sonarqube impl.Sonarqube
	dataflowTester := e2ehelper.NewDataFlowTester(t, "sonarqube", sonarqube)
	basicRes := implcontext.NewDefaultBasicRes(dataflowTester.Cfg, dataflowTester.Log, dataflowTester.Dal)
	const name = "modify _tool_sonarqube_file_metrics file_metrics_key from varchar(191) to varchar(500)"

	// normal run: the primary key exists and is re-created
	prepareFileMetricsBeforeKeyLength(t, dataflowTester)
	runMigration(t, sonarqubemigrations.All(), name, basicRes)
	assertFileMetricsMigrated(t, dataflowTester)

	// re-run: an interrupted attempt dropped the primary key before the
	// migration history was recorded
	prepareFileMetricsBeforeKeyLength(t, dataflowTester)
	if dataflowTester.Dal.Dialect() == "mysql" {
		require.NoError(t, dataflowTester.Dal.Exec("ALTER TABLE _tool_sonarqube_file_metrics DROP PRIMARY KEY"))
	} else {
		require.NoError(t, dataflowTester.Dal.Exec("ALTER TABLE _tool_sonarqube_file_metrics DROP CONSTRAINT _tool_sonarqube_file_metrics_pkey"))
	}
	require.False(t, hasPrimaryKey(t, dataflowTester))
	runMigration(t, sonarqubemigrations.All(), name, basicRes)
	assertFileMetricsMigrated(t, dataflowTester)
}

func prepareFileMetricsBeforeKeyLength(t *testing.T, dataflowTester *e2ehelper.DataFlowTester) {
	t.Helper()
	require.NoError(t, dataflowTester.Db.Migrator().DropTable(&sonarqubeFileMetricsBeforeKeyLength{}))
	require.NoError(t, dataflowTester.Db.AutoMigrate(&sonarqubeFileMetricsBeforeKeyLength{}))
	require.NoError(t, dataflowTester.Db.Create(&sonarqubeFileMetricsBeforeKeyLength{
		ConnectionId:   1,
		FileMetricsKey: "project:src/main.go",
		ProjectKey:     "project",
	}).Error)
}

func assertFileMetricsMigrated(t *testing.T, dataflowTester *e2ehelper.DataFlowTester) {
	t.Helper()
	require.True(t, hasPrimaryKey(t, dataflowTester))
	columnTypes, err := dataflowTester.Db.Migrator().ColumnTypes(&sonarqubeFileMetricsBeforeKeyLength{})
	require.NoError(t, err)
	for _, columnType := range columnTypes {
		if columnType.Name() == "file_metrics_key" {
			length, ok := columnType.Length()
			require.True(t, ok)
			require.Equal(t, int64(500), length)
		}
	}
	var migrated sonarqubeFileMetricsBeforeKeyLength
	require.NoError(t, dataflowTester.Db.First(&migrated, "connection_id = ? AND file_metrics_key = ?", 1, "project:src/main.go").Error)
	require.Equal(t, "project", migrated.ProjectKey)
}

func hasPrimaryKey(t *testing.T, dataflowTester *e2ehelper.DataFlowTester) bool {
	t.Helper()
	schema := "current_schema()"
	if dataflowTester.Dal.Dialect() == "mysql" {
		schema = "DATABASE()"
	}
	var count int64
	require.NoError(t, dataflowTester.Db.Raw(
		"SELECT COUNT(*) FROM information_schema.table_constraints WHERE table_schema = "+schema+
			" AND table_name = ? AND constraint_type = 'PRIMARY KEY'",
		"_tool_sonarqube_file_metrics",
	).Scan(&count).Error)
	return count > 0
}
