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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/core"
)

type changeIndexOfJobPath struct{}

type jenkinsJob20221110Before struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	FullName     string `gorm:"primaryKey;type:varchar(255)"`
	Name         string `gorm:"index;type:varchar(255)"`
	Path         string `gorm:"primaryKey;type:varchar(511)"`
	Class        string `gorm:"type:varchar(255)"`
	Color        string `gorm:"type:varchar(255)"`
	Base         string `gorm:"type:varchar(255)"`
	Url          string
	Description  string
	PrimaryView  string `gorm:"type:varchar(255)"`
	archived.NoPKModel
}

func (jenkinsJob20221110Before) TableName() string {
	return "_tool_jenkins_jobs"
}

type jenkinsJob20221110After struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	FullName     string `gorm:"primaryKey;type:varchar(255)"`
	Name         string `gorm:"index;type:varchar(255)"`
	Path         string `gorm:"index;type:varchar(511)"`
	Class        string `gorm:"type:varchar(255)"`
	Color        string `gorm:"type:varchar(255)"`
	Base         string `gorm:"type:varchar(255)"`
	Url          string
	Description  string
	PrimaryView  string `gorm:"type:varchar(255)"`
	archived.NoPKModel
}

func (jenkinsJob20221110After) TableName() string {
	return "_tool_jenkins_jobs"
}

func (script *changeIndexOfJobPath) Up(basicRes core.BasicRes) errors.Error {
	return migrationhelper.TransformTable(
		basicRes,
		script,
		"_tool_jenkins_jobs",
		func(s *jenkinsJob20221110Before) (*jenkinsJob20221110After, errors.Error) {
			dst := &jenkinsJob20221110After{
				ConnectionId: s.ConnectionId,
				FullName:     s.FullName,
				Name:         s.Name,
				Path:         s.Path,
				Class:        s.Class,
				Color:        s.Color,
				Base:         s.Base,
				Url:          s.Url,
				Description:  s.Description,
				PrimaryView:  s.PrimaryView,
				NoPKModel:    s.NoPKModel,
			}
			return dst, nil
		},
	)
}

func (*changeIndexOfJobPath) Version() uint64 {
	return 20221110231237
}

func (*changeIndexOfJobPath) Name() string {
	return "add url to jenkinsJob"
}
