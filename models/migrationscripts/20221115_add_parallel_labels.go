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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"time"
)

type DbPipelineParallelLabel20221115 struct {
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	PipelineId uint64    `json:"pipeline_id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"primaryKey"`
}

func (DbPipelineParallelLabel20221115) TableName() string {
	return "_devlake_pipeline_parallel_labels"
}

type DbBlueprintParallelLabel20221115 struct {
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	BlueprintId uint64    `json:"blueprint_id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"primaryKey"`
}

func (DbBlueprintParallelLabel20221115) TableName() string {
	return "_devlake_blueprint_parallel_labels"
}

type addParallelLabels struct{}

func (*addParallelLabels) Up(res core.BasicRes) errors.Error {
	db := res.GetDal()
	err := db.AutoMigrate(&DbPipelineParallelLabel20221115{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&DbBlueprintParallelLabel20221115{})
	if err != nil {
		return err
	}
	return nil
}

func (*addParallelLabels) Version() uint64 {
	return 20221115000034
}

func (*addParallelLabels) Name() string {
	return "UpdateSchemas for addParallelLabels"
}
