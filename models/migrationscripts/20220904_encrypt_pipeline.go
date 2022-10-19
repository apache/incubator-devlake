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
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/core"
	"gorm.io/datatypes"
)

var _ core.MigrationScript = (*encryptPipeline)(nil)

type encryptPipeline struct{}

type Pipeline20220904Before struct {
	archived.Model
	Name          string         `json:"name" gorm:"index"`
	BlueprintId   uint64         `json:"blueprintId"`
	Plan          datatypes.JSON `json:"plan"` // target field
	TotalTasks    int            `json:"totalTasks"`
	FinishedTasks int            `json:"finishedTasks"`
	BeganAt       *time.Time     `json:"beganAt"`
	FinishedAt    *time.Time     `json:"finishedAt" gorm:"index"`
	Status        string         `json:"status"`
	Message       string         `json:"message"`
	SpentSeconds  int            `json:"spentSeconds"`
	Stage         int            `json:"stage"`
}

func (Pipeline20220904Before) TableName() string {
	return "_devlake_pipelines"
}

type Pipeline0904After struct {
	common.Model
	Name          string     `json:"name" gorm:"index"`
	BlueprintId   uint64     `json:"blueprintId"`
	Plan          string     `json:"plan" encrypt:"yes"` // target field
	TotalTasks    int        `json:"totalTasks"`
	FinishedTasks int        `json:"finishedTasks"`
	BeganAt       *time.Time `json:"beganAt"`
	FinishedAt    *time.Time `json:"finishedAt" gorm:"index"`
	Status        string     `json:"status"`
	Message       string     `json:"message"`
	SpentSeconds  int        `json:"spentSeconds"`
	Stage         int        `json:"stage"`
}

func (Pipeline0904After) TableName() string {
	return "_devlake_pipelines"
}

func (script *encryptPipeline) Up(basicRes core.BasicRes) errors.Error {
	encKey := basicRes.GetConfig(core.EncodeKeyEnvStr)
	if encKey == "" {
		return errors.BadInput.New("invalid encKey")
	}

	return migrationhelper.TransformTable(
		basicRes,
		script,
		"_devlake_pipelines",
		func(s *Pipeline20220904Before) (*Pipeline0904After, errors.Error) {
			encryptedPlan, err := core.Encrypt(encKey, string(s.Plan))
			if err != nil {
				return nil, err
			}

			dst := &Pipeline0904After{
				Name:          s.Name,
				BlueprintId:   s.BlueprintId,
				FinishedTasks: s.FinishedTasks,
				BeganAt:       s.BeganAt,
				FinishedAt:    s.FinishedAt,
				Status:        s.Status,
				Message:       s.Message,
				SpentSeconds:  s.SpentSeconds,
				Stage:         s.Stage,
				Plan:          encryptedPlan,
			}
			return dst, nil
		},
	)

}

func (*encryptPipeline) Version() uint64 {
	return 20220904162121
}

func (*encryptPipeline) Name() string {
	return "encrypt Pipeline"
}
