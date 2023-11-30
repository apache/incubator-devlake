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
	"time"
)

type circleciJob20231129 struct {
	CreatedAt *time.Time `json:"created_at"`
	QueuedAt  *time.Time `json:"queued_at"`
	Duration  int64      `json:"duration"`
}

func (circleciJob20231129) TableName() string {
	return "_tool_circleci_jobs"
}

type addFieldsToCircleciJob20231129 struct{}

func (*addFieldsToCircleciJob20231129) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(&circleciJob20231129{})
}

func (*addFieldsToCircleciJob20231129) Version() uint64 {
	return 20231129211000
}

func (*addFieldsToCircleciJob20231129) Name() string {
	return "add some fields to _tool_circleci_jobs"
}
