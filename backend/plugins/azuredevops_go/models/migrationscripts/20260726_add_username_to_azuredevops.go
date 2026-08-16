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

type addUsernameToAzuredevops struct{}

type azuredevopsConnection20260726 struct {
	Username string `gorm:"type:varchar(255)"`
}

func (azuredevopsConnection20260726) TableName() string {
	return "_tool_azuredevops_go_connections"
}

func (script *addUsernameToAzuredevops) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&azuredevopsConnection20260726{},
	)
}

func (*addUsernameToAzuredevops) Version() uint64 {
	return 20260726000001
}

func (*addUsernameToAzuredevops) Name() string {
	return "add username field to _tool_azuredevops_go_connections"
}
