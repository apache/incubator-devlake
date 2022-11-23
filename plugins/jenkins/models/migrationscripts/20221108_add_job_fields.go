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
)

type addJobFields struct{}

type jenkinsJob20221108 struct {
	Name        string `gorm:"index;type:varchar(255)"`
	Url         string
	Description string
	PrimaryView string `gorm:"type:varchar(255)"`
}

func (jenkinsJob20221108) TableName() string {
	return "_tool_jenkins_jobs"
}

func (script *addJobFields) Up(basicRes core.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.RenameColumn(`_tool_jenkins_jobs`, `name`, `full_name`)
	if err != nil {
		return err
	}
	return db.AutoMigrate(&jenkinsJob20221108{})
}

func (*addJobFields) Version() uint64 {
	return 20221108231237
}

func (*addJobFields) Name() string {
	return "add fields for jenkins jobs"
}
