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
	"context"
	"github.com/apache/incubator-devlake/errors"

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type addPipeline20220914 struct{}

type BitbucketPipeline20220914 struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	PipelineId   int    `gorm:"primaryKey"`
	ProjectId    int    `gorm:"primaryKey;type:varchar(255)"`
	Ref          string `gorm:"type:varchar(255)"`
	Sha          string `gorm:"type:varchar(255)"`
	archived.NoPKModel
}

func (BitbucketPipeline20220914) TableName() string {
	return "_tool_bitbucket_pipelines"
}

func (*addPipeline20220914) Up(ctx context.Context, db *gorm.DB) errors.Error {
	err := db.Migrator().CreateTable(&BitbucketPipeline20220914{})
	if err != nil {
		return errors.Convert(err)
	}
	return nil
}

func (*addPipeline20220914) Version() uint64 {
	return 20220914111223
}

func (*addPipeline20220914) Name() string {
	return "bitbucket add _tool_bitbucket_pipelines table"
}
