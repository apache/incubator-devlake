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
)

type addTaskLeft struct{}

type ZentaoTask20230627 struct {
	Left float64 `json:"left" gorm:"column:db_left"`
}

func (ZentaoTask20230627) TableName() string {
	return "_tool_zentao_tasks"
}

func (*addTaskLeft) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&ZentaoTask20230627{},
	)
}

func (*addTaskLeft) Version() uint64 {
	return 20230627000001
}

func (*addTaskLeft) Name() string {
	return "zentao init changelog schemas"
}
