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
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/plugins/core"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/common"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Pipeline0904TempBefore struct {
	NewPlan string
}

func (Pipeline0904TempBefore) TableName() string {
	return "_devlake_pipelines"
}

type Pipeline0904Temp struct {
	common.Model
	Name          string     `json:"name" gorm:"index"`
	BlueprintId   uint64     `json:"blueprintId"`
	NewPlan       string     `json:"plan" encrypt:"yes"`
	TotalTasks    int        `json:"totalTasks"`
	FinishedTasks int        `json:"finishedTasks"`
	BeganAt       *time.Time `json:"beganAt"`
	FinishedAt    *time.Time `json:"finishedAt" gorm:"index"`
	Status        string     `json:"status"`
	Message       string     `json:"message"`
	SpentSeconds  int        `json:"spentSeconds"`
	Stage         int        `json:"stage"`
	Plan          datatypes.JSON
}

func (Pipeline0904Temp) TableName() string {
	return "_devlake_pipelines"
}

type Pipeline0904TempAfter struct {
	NewPlan string
	Plan    string
}

func (Pipeline0904TempAfter) TableName() string {
	return "_devlake_pipelines"
}

type encryptPipeline struct{}

func (*encryptPipeline) Up(ctx context.Context, db *gorm.DB) errors.Error {
	err := db.AutoMigrate(&Pipeline0904TempBefore{})
	if err != nil {
		return errors.Convert(err)
	}
	var result *gorm.DB
	var pipelineList []Pipeline0904Temp
	result = db.Find(&pipelineList)

	if result.Error != nil {
		return errors.Convert(result.Error)
	}

	// Encrypt all pipelines.plan&settings which had been stored before v0.14
	for _, v := range pipelineList {
		c := config.GetConfig()
		encKey := c.GetString(core.EncodeKeyEnvStr)
		if encKey == "" {
			return errors.BadInput.New("invalid encKey")
		}
		v.NewPlan, err = core.Encrypt(encKey, string(v.Plan))
		if err != nil {
			return errors.Convert(err)
		}

		err = errors.Convert(db.Save(&v).Error)
		if err != nil {
			return errors.Convert(err)
		}
	}

	err = db.Migrator().DropColumn(&Pipeline0904Temp{}, "plan")
	if err != nil {
		return errors.Convert(err)
	}
	err = db.Migrator().RenameColumn(&Pipeline0904TempAfter{}, "new_plan", "plan")
	if err != nil {
		return errors.Convert(err)
	}
	_ = db.Find(&pipelineList)
	return nil
}

func (*encryptPipeline) Version() uint64 {
	return 20220904162121
}

func (*encryptPipeline) Name() string {
	return "encrypt Pipeline"
}
