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
	"time"
)

var _ plugin.MigrationScript = (*addQueuedFieldsInJobBuild20231128)(nil)

type addQueuedFieldsInJobBuild20231128 struct{}

func (*addQueuedFieldsInJobBuild20231128) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(&bambooJobBuild20231128{})
}

func (*addQueuedFieldsInJobBuild20231128) Version() uint64 {
	return 20231128202200
}

func (*addQueuedFieldsInJobBuild20231128) Name() string {
	return "add queue fields in _tool_bamboo_job_builds"
}

type bambooJobBuild20231128 struct {
	QueueStartedTime       *time.Time `json:"queueStartedTime"`
	QueueDuration          int64      `json:"queueDuration"`
	QueueDurationInSeconds int64      `json:"queueDurationInSeconds"`
}

func (bambooJobBuild20231128) TableName() string {
	return "_tool_bamboo_job_builds"
}
