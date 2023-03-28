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

package tasks

import (
	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestBuildDeploymentPairs(t *testing.T) {
	db := new(mockdal.Dal)
	data := &DoraTaskData{
		Options: &DoraOptions{
			ProjectName: "project1",
		},
	}
	deploymentDiffPairs := []deploymentPair{
		{TaskId: "1", RepoId: "repo1", NewDeployCommitSha: "sha1", OldDeployCommitSha: ""},
		{TaskId: "2", RepoId: "repo1", NewDeployCommitSha: "sha2", OldDeployCommitSha: "sha1"},
		{TaskId: "3", RepoId: "repo2", NewDeployCommitSha: "sha3", OldDeployCommitSha: ""},
		{TaskId: "4", RepoId: "repo2", NewDeployCommitSha: "sha4", OldDeployCommitSha: "sha3"},
	}

	expectedPairs := []deploymentPair{
		{TaskId: "1", RepoId: "repo1", NewDeployCommitSha: "sha1", OldDeployCommitSha: ""},
		{TaskId: "2", RepoId: "repo1", NewDeployCommitSha: "sha2", OldDeployCommitSha: "sha1"},
		{TaskId: "3", RepoId: "repo2", NewDeployCommitSha: "sha3", OldDeployCommitSha: ""},
		{TaskId: "4", RepoId: "repo2", NewDeployCommitSha: "sha4", OldDeployCommitSha: "sha3"},
	}
	db.On("All", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		dst := args.Get(0).(*[]deploymentPair)
		*dst = deploymentDiffPairs
	}).Return(nil).Once()
	res, err := buildDeploymentPairs(db, data)

	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(res, expectedPairs))

	if !reflect.DeepEqual(deploymentDiffPairs, expectedPairs) {
		t.Errorf("buildDeploymentPairs() = %v, want %v", deploymentDiffPairs, expectedPairs)
	}
}
