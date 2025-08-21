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
	"encoding/json"
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*addStore)(nil)

type store20240318 struct {
	StoreKey   string          `gorm:"primaryKey;type:varchar(255)"`
	StoreValue json.RawMessage `gorm:"type:json;serializer:json"`
	CreatedAt  time.Time       `json:"createdAt" mapstructure:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt" mapstructure:"updatedAt"`
}

func (store20240318) TableName() string {
	return "_devlake_store"
}

type addStore struct{}

func (*addStore) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(store20240318{})
}

func (*addStore) Version() uint64 {
	return 20240318111247
}

func (*addStore) Name() string {
	return "add _devlake_store table"
}
