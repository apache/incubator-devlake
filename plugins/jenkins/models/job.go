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
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/helper"
)

// JenkinsJobProps current used jenkins job props
type JenkinsJobProps struct {
	// collected fields
	ConnectionId        uint64 `gorm:"primaryKey"`
	Name                string `gorm:"primaryKey;type:varchar(255)"`
	Path                string `gorm:"primaryKey;type:varchar(511)"`
	Class               string `gorm:"type:varchar(255)"`
	Color               string `gorm:"type:varchar(255)"`
	Base                string `gorm:"type:varchar(255)"`
	HasUpstreamProjects bool
}

// JenkinsJob db entity for jenkins job
type JenkinsJob struct {
	JenkinsJobProps
	common.NoPKModel
}

func (JenkinsJob) TableName() string {
	return "_tool_jenkins_jobs"
}

type FolderInput struct {
	*helper.ListBaseNode
	Path string
}

func (f *FolderInput) Data() interface{} {
	return f.Path
}

func NewFolderInput(path string) *FolderInput {
	return &FolderInput{
		Path:         path,
		ListBaseNode: helper.NewListBaseNode(),
	}
}
