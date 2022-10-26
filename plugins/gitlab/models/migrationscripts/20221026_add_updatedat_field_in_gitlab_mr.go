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

type mergeRequest20221026 struct {
	GitlabUpdatedAt time.Time
}

func (mergeRequest20221026) TableName() string {
	return "_tool_gitlab_merge_requests"
}

type addGitlabUpdateAtInGitlabMr struct{}

func (*addGitlabUpdateAtInGitlabMr) Up(basicRes core.BasicRes) errors.Error {
	err := basicRes.GetDal().AutoMigrate(&mergeRequest20221026{})
	if err != nil {
		return err
	}
	return nil
}

func (*addGitlabUpdateAtInGitlabMr) Version() uint64 {
	return 20221026145320
}

func (*addGitlabUpdateAtInGitlabMr) Name() string {
	return "add column `gitlab_updated_at` at _tool_gitlab_merge_requests"
}
