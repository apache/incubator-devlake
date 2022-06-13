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

package archived

import "github.com/apache/incubator-devlake/models/migrationscripts/archived"

// This Model is intended to save commits that are associated to a merge request
// for the purpose of calculating PR Review Rounds and other metrics that
// rely on associating commits to merge requests that may or may not
// exist on the main branch of a project.
// Thus a "Merge Request Commit" needs to be considered as distinct from a "Commit"

type GitlabMergeRequestCommit struct {
	ConnectionId   uint64 `gorm:"primaryKey"`
	CommitSha      string `gorm:"primaryKey;type:varchar(40)"`
	MergeRequestId int    `gorm:"primaryKey;autoIncrement:false"`
	archived.NoPKModel
}

func (GitlabMergeRequestCommit) TableName() string {
	return "_tool_gitlab_merge_request_commits"
}
