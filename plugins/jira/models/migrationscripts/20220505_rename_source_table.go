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
	"github.com/apache/incubator-devlake/plugins/core"
)

type renameSourceTable20220505 struct{}

func (*renameSourceTable20220505) Up(basicRes core.BasicRes) errors.Error {
	return basicRes.GetDal().RenameTable("_tool_jira_sources", "_tool_jira_connections")
}

func (*renameSourceTable20220505) Version() uint64 {
	return 20220505212344
}

func (*renameSourceTable20220505) Name() string {
	return "Rename source to connection "
}
