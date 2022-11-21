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

type DbPipelineLabel20221115 struct {
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	PipelineId uint64    `json:"pipeline_id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"primaryKey"`
}

func (DbPipelineLabel20221115) TableName() string {
	return "_devlake_pipeline_labels"
}

type DbBlueprintLabel20221115 struct {
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	BlueprintId uint64    `json:"blueprint_id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"primaryKey"`
}

func (DbBlueprintLabel20221115) TableName() string {
	return "_devlake_blueprint_labels"
}

type addLabels struct{}

func (*addLabels) Up(res core.BasicRes) errors.Error {
	db := res.GetDal()
	err := db.AutoMigrate(&DbPipelineLabel20221115{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&DbBlueprintLabel20221115{})
	if err != nil {
		return err
	}
	return nil
}

func (*addLabels) Version() uint64 {
	return 20221115000034
}

func (*addLabels) Name() string {
	return "add parallel labels' schema for blueprint and pipeline"
}
