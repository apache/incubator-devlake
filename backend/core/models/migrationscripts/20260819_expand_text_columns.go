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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*expandDomainTextColumns)(nil)

type expandDomainTextColumns struct{}

func (*expandDomainTextColumns) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	columns := []struct {
		tableName  string
		columnName string
	}{
		{"cicd_tasks", "name"},
		{"cicd_scopes", "name"},
		{"cicd_releases", "name"},
		{"cicd_releases", "display_title"},
		{"cicd_deployment_commits", "name"},
		{"cicd_deployment_commits", "subtask_name"},
		{"cicd_deployment_commits", "ref_name"},
		{"incidents", "component"},
	}
	for _, column := range columns {
		if err := db.ModifyColumnType(column.tableName, column.columnName, "text"); err != nil {
			return err
		}
	}
	return nil
}

func (*expandDomainTextColumns) Version() uint64 {
	return 20260819000001
}

func (*expandDomainTextColumns) Name() string {
	return "expand domain text columns"
}
