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

package models

import (
	"sync"
	"testing"

	"github.com/apache/incubator-devlake/plugins/kiro/models/migrationscripts/archived"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mysqlDriver "gorm.io/driver/mysql"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm/schema"
)

func TestLogTimestampUsesPortableMicrosecondPrecision(t *testing.T) {
	logModels := []struct {
		name  string
		model any
	}{
		{name: "chat runtime model", model: &KiroChatLog{}},
		{name: "completion runtime model", model: &KiroCompletionLog{}},
		{name: "chat migration model", model: &archived.KiroChatLog{}},
		{name: "completion migration model", model: &archived.KiroCompletionLog{}},
	}

	for _, logModel := range logModels {
		t.Run(logModel.name, func(t *testing.T) {
			sch, err := schema.Parse(logModel.model, &sync.Map{}, schema.NamingStrategy{})
			require.NoError(t, err)

			field := sch.LookUpField("Timestamp")
			require.NotNil(t, field)
			assert.Equal(t, 6, field.Precision)
			assert.Equal(t, "datetime(6) NULL", mysqlDriver.New(mysqlDriver.Config{}).DataTypeOf(field))
			assert.Equal(t, "timestamptz(6)", postgresDriver.New(postgresDriver.Config{}).DataTypeOf(field))
		})
	}
}
