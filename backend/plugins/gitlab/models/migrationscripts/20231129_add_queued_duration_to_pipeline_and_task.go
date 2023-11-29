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

var _ plugin.MigrationScript = (*addQueuedDuration20231129)(nil)

type gitlabJob20231129 struct {
	QueuedDuration float64
}

func (gitlabJob20231129) TableName() string {
	return "_tool_gitlab_jobs"
}

type gitlabPipeline20231129 struct {
	QueuedDuration float64
}

func (gitlabPipeline20231129) TableName() string {
	return "_tool_gitlab_pipelines"
}

type addQueuedDuration20231129 struct{}

func (script *addQueuedDuration20231129) Up(basicRes context.BasicRes) errors.Error {
	if err := basicRes.GetDal().AutoMigrate(&gitlabJob20231129{}); err != nil {
		return err
	}
	if err := basicRes.GetDal().AutoMigrate(&gitlabPipeline20231129{}); err != nil {
		return err
	}
	return nil
}

func (*addQueuedDuration20231129) Version() uint64 {
	return 20231129110000
}

func (script *addQueuedDuration20231129) Name() string {
	return "add queued_duration field to _tool_gitlab_jobs and _tool_gitlab_pipelines"
}
