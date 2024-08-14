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

var _ plugin.MigrationScript = (*changeIssueComponentType)(nil)

type changeIssueComponentType struct{}

func (script *changeIssueComponentType) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().ModifyColumnType("_tool_github_issues", "components", "text")
}

func (*changeIssueComponentType) Version() uint64 {
	return 20240813154633
}

func (*changeIssueComponentType) Name() string {
	return "change _tool_github_issues.components type to text"
}
