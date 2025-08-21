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

var _ plugin.MigrationScript = (*removeMrEnricherFields)(nil)

type removeMrEnricherFields struct{}

func (script *removeMrEnricherFields) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().DropColumns("_tool_gitlab_merge_requests", "first_comment_time", "review_rounds")
}

func (*removeMrEnricherFields) Version() uint64 {
	return 20240904152414
}

func (*removeMrEnricherFields) Name() string {
	return "remove mr review fields in favor of dora plugin"
}
