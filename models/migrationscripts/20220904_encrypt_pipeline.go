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

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/common"
	"gorm.io/gorm"
)

type PipelinEncryption0904 struct {
	common.Model
	Plan    string
	RawPlan string
}

func (PipelinEncryption0904) TableName() string {
	return "_devlake_pipelines"
}

type encryptPipeline struct{}

func (*encryptPipeline) Up(ctx context.Context, db *gorm.DB) errors.Error {
	c := config.GetConfig()
	encKey := c.GetString(core.EncodeKeyEnvStr)
	if encKey == "" {
		return errors.BadInput.New("invalid encKey")
	}

	pipeline := &PipelinEncryption0904{}
	err := db.Migrator().RenameColumn(pipeline, "plan", "raw_plan")
	if err != nil {
		return errors.Convert(err)
	}
	err = db.AutoMigrate(pipeline)
	if err != nil {
		return errors.Convert(err)
	}

	// Encrypt all pipelines.plan which had been stored before v0.14
	cursor, err := db.Model(pipeline).Rows()
	if err != nil {
		return errors.Convert(err)
	}
	defer cursor.Close()

	for cursor.Next() {
		err = db.ScanRows(cursor, pipeline)
		if err != nil {
			return errors.Convert(err)
		}
		pipeline.Plan, err = core.Encrypt(encKey, string(pipeline.RawPlan))
		if err != nil {
			return errors.Convert(err)
		}
		err = errors.Convert(db.Save(pipeline).Error)
		if err != nil {
			return errors.Convert(err)
		}
	}

	err = db.Migrator().DropColumn(pipeline, "raw_plan")
	if err != nil {
		return errors.Convert(err)
	}
	_ = db.First(pipeline)
	return nil
}

func (*encryptPipeline) Version() uint64 {
	return 20220904162121
}

func (*encryptPipeline) Name() string {
	return "encrypt Pipeline"
}
