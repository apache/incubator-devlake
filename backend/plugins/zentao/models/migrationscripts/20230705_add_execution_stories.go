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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/zentao/models/archived"
)

type addExecutionStoryAndExecutionSummary struct{}

func (*addExecutionStoryAndExecutionSummary) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&archived.ZentaoExecutionStory{},
		&archived.ZentaoExecutionSummary{},
		&archived.ZentaoProductSummary{},
		&archived.ZentaoProjectStory{},
	)
}

func (*addExecutionStoryAndExecutionSummary) Version() uint64 {
	return 20230705135809
}

func (*addExecutionStoryAndExecutionSummary) Name() string {
	return "add table _tool_zentao_execution_stories, _tool_zentao_execution_summary and _tool_zentao_product_summary"
}
