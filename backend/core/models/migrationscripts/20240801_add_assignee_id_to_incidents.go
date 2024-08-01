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

var _ plugin.MigrationScript = (*addAssigneeToIncident)(nil)

type incident20240801 struct {
	AssigneeId   string `gorm:"type:varchar(255)"`
	AssigneeName string `gorm:"type:varchar(255)"`
}

func (incident20240801) TableName() string {
	return "incidents"
}

type addAssigneeToIncident struct{}

func (*addAssigneeToIncident) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(&incident20240801{})
}

func (*addAssigneeToIncident) Version() uint64 {
	return 20240801110000
}

func (*addAssigneeToIncident) Name() string {
	return "add assignee info to incidents"
}
