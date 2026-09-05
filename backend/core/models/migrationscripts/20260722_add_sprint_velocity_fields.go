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
	"github.com/apache/devlake/core/context"
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
)

var _ plugin.MigrationScript = (*addSprintVelocityFields)(nil)

type sprint20260722 struct {
	CommittedStoryPoint *float64
	CompletedStoryPoint *float64
}

func (sprint20260722) TableName() string {
	return "sprints"
}

type addSprintVelocityFields struct{}

func (script *addSprintVelocityFields) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(new(sprint20260722))
}

func (*addSprintVelocityFields) Version() uint64 {
	return 20260722100000
}

func (*addSprintVelocityFields) Name() string {
	return "add committed_story_point/completed_story_point to sprints"
}
