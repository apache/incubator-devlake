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

var _ plugin.MigrationScript = (*modifyFileMetricsKeyLength)(nil)

type modifyFileMetricsKeyLength struct{}

type fileMetricsKey20230927 struct {
	FileMetricsKey string `gorm:"primaryKey;type:varchar(500)"`
}

func (fileMetricsKey20230927) TableName() string {
	return "_tool_sonarqube_file_metrics"
}

func (script *modifyFileMetricsKeyLength) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	tx := db.Begin()
	err := migrationhelper.ChangeColumnsType[fileMetricsKey20230927](
		basicRes,
		script,
		fileMetricsKey20230927{}.TableName(),
		[]string{"file_metrics_key"},
		func(tmpColumnParams []interface{}) errors.Error {
			return db.UpdateColumn(
				&fileMetricsKey20230927{},
				"file_metrics_key",
				dal.DalClause{Expr: " ? ", Params: tmpColumnParams},
				dal.Where("? != '' ", tmpColumnParams...),
			)
		},
	)
	if err != nil {
		return tx.Rollback()
	}

	if err = tx.Exec("ALTER TABLE _tool_sonarqube_file_metrics DROP PRIMARY KEY"); err != nil {
		return tx.Rollback()
	}

	if err := tx.Exec("ALTER TABLE _tool_sonarqube_file_metrics ADD PRIMARY KEY (connection_id, file_metrics_key)"); err != nil {
		return tx.Rollback()
	}

	if err := tx.Commit(); err != nil {
		return tx.Rollback()
	}

	return nil
}

func (*modifyFileMetricsKeyLength) Version() uint64 {
	return 20230927145127
}

func (*modifyFileMetricsKeyLength) Name() string {
	return "modify _tool_sonarqube_file_metrics file_metrics_key from varchar(191) to varchar(500)"
}
