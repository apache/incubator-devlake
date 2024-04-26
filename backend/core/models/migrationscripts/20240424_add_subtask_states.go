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

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*addSubtaskStates)(nil)

type subtaskState20240424 struct {
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	Plugin             string    `gorm:"primaryKey;type:varchar(50)" json:"plugin"`
	Subtask            string    `gorm:"primaryKey;type:varchar(50)" json:"subtask"`
	RawDataTable       string    `gorm:"primaryKey;column:raw_data_table;type:varchar(50)" json:"raw_data_table"`
	RawDataParams      string    `gorm:"primaryKey;column:raw_data_params;type:varchar(255);index" json:"raw_data_params"`
	TimeAfter          *time.Time
	LatestSuccessStart *time.Time
}

func (subtaskState20240424) TableName() string {
	return "_devlake_subtask_states"
}

type addSubtaskStates struct{}

func (script *addSubtaskStates) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	errors.Must(db.AutoMigrate(&subtaskState20240424{}))
	return nil
}

func (*addSubtaskStates) Version() uint64 {
	return 20240424152734
}

func (*addSubtaskStates) Name() string {
	return "add _devlake_subtask_states"
}
