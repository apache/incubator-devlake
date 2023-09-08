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

type addFullSyncToBlueprint struct {
	FullSync bool `json:"fullSync"`
}

func (*addFullSyncToBlueprint) TableName() string {
	return "_devlake_blueprints"
}

type addFullSyncToPipeline struct {
	FullSync bool `json:"fullSync"`
}

func (*addFullSyncToPipeline) TableName() string {
	return "_devlake_pipelines"
}

type addFullSync struct{}

func (*addFullSync) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&addFullSyncToBlueprint{},
		&addFullSyncToPipeline{},
	)
}

func (*addFullSync) Version() uint64 {
	return 20230907000041
}

func (*addFullSync) Name() string {
	return "add full_sync to _devlake_blueprints and _devlake_pipelines table"
}
