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
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
)

type workflow20240717 struct {
	CreatedDate *time.Time
}

func (workflow20240717) TableName() string {
	return "_tool_circleci_workflows"
}

type addCreatedStoppedDateToWorkflow struct{}

func (*addCreatedStoppedDateToWorkflow) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if err := db.AutoMigrate(&workflow20240717{}); err != nil {
		return err
	}
	return db.RenameColumn(workflow20240717{}.TableName(), "stopped_at", "stopped_date")
}

func (*addCreatedStoppedDateToWorkflow) Version() uint64 {
	return 20240717210714
}

func (*addCreatedStoppedDateToWorkflow) Name() string {
	return "add created_date and stopped_date to _tool_circleci_workflows, drop stopped_at column"
}
