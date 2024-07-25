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

var _ plugin.MigrationScript = (*addMergedByToPr)(nil)

type pr20240710 struct {
	MergedByName string `gorm:"type:varchar(100)"`
	MergedById   int
}

func (pr20240710) TableName() string {
	return "_tool_github_pull_requests"
}

type addMergedByToPr struct{}

func (*addMergedByToPr) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if err := db.AutoMigrate(&pr20240710{}); err != nil {
		return err
	}
	return nil
}

func (*addMergedByToPr) Version() uint64 {
	return 20240710142100
}

func (*addMergedByToPr) Name() string {
	return "add merged by to _tool_github_pull_requests"
}
