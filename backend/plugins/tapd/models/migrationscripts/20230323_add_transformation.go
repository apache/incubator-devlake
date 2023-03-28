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
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

type TapdWorkspace20230323 struct {
	TransformationRuleId uint64 `json:"transformationRuleId,omitempty" mapstructure:"transformationRuleId"`
}

func (TapdWorkspace20230323) TableName() string {
	return "_tool_tapd_workspaces"
}

type addTransformation struct{}

func (*addTransformation) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&models.TapdTransformationRule{},
		&TapdWorkspace20230323{},
	)
}

func (*addTransformation) Version() uint64 {
	return 20230323000003
}

func (*addTransformation) Name() string {
	return "Tapd add transformation"
}
