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

type fixRunNameToText struct{}

type githubRun20230518_old struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	RepoId       int    `gorm:"primaryKey"`
	ID           int    `gorm:"primaryKey;autoIncrement:false"`
	Name         string `gorm:"type:varchar(255)"`
}
type githubRun20230518 struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	RepoId       int    `gorm:"primaryKey"`
	ID           int    `gorm:"primaryKey;autoIncrement:false"`
	Name         string `gorm:"type:text"`
}

func (*fixRunNameToText) Up(baseRes context.BasicRes) errors.Error {
	err := migrationhelper.TransformColumns(
		baseRes,
		&fixRunNameToText{},
		"_tool_github_runs",
		[]string{"name"},
		func(src *githubRun20230518_old) (*githubRun20230518, errors.Error) {
			return &githubRun20230518{
				ConnectionId: src.ConnectionId,
				RepoId:       src.RepoId,
				ID:           src.ID,
				Name:         src.Name,
			}, nil
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (*fixRunNameToText) Version() uint64 {
	return 20230518000002
}

func (*fixRunNameToText) Name() string {
	return "UpdateSchemas for fixRunNameToText"
}
