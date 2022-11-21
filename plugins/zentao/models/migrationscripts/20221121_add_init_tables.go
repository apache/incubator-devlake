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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/zentao/models/archived"
)

type addInitTables struct{}

func (*addInitTables) Up(basicRes core.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.DropTables(
		&archived.ZentaoConnection{},
		&archived.ZentaoProject{},
		&archived.ZentaoProduct{},
		&archived.ZentaoExecution{},
		&archived.ZentaoStory{},
		&archived.ZentaoBug{},
		&archived.ZentaoTask{},
		"_tool_zentao_bugs`",
		"_tool_zentao_executions`",
		"_tool_zentao_products`",
		"_tool_zentao_stories`",
		"_tool_zentao_tasks`",
	)
	if err != nil {
		return err
	}
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&archived.ZentaoConnection{},
		&archived.ZentaoProject{},
		&archived.ZentaoProduct{},
		&archived.ZentaoExecution{},
		&archived.ZentaoStory{},
		&archived.ZentaoBug{},
		&archived.ZentaoTask{},
		&archived.ZentaoAccount{},
		&archived.ZentaoDepartment{},
	)
}

func (*addInitTables) Version() uint64 {
	return 20221121000001
}

func (*addInitTables) Name() string {
	return "zentao init schemas"
}
