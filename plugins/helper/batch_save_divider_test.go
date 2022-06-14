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

package helper

import (
	"reflect"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/helpers/unithelper"
	"github.com/apache/incubator-devlake/mocks"
	"github.com/apache/incubator-devlake/models/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockJirIssueBsd struct {
	common.RawDataOrigin
	ID               string `gorm:"primaryKey"`
	Title            string
	ChangelogUpdated *time.Time
}

type MockJiraChangelogBsd struct {
	common.RawDataOrigin
	ID string `gorm:"primaryKey"`
}

type MockJiraIssueChangelogUpdatedBsd struct {
	ID               string `gorm:"primaryKey"`
	ChangelogUpdated *time.Time
}

func (MockJiraIssueChangelogUpdatedBsd) PartialUpdate() {}

func TestBatchSaveDivider(t *testing.T) {
	mockDal := new(mocks.Dal)

	mockLog := unithelper.DummyLogger()
	mockRes := new(mocks.BasicRes)
	mockRes.On("GetDal").Return(mockDal)
	mockRes.On("GetLogger").Return(mockLog)

	divider := NewBatchSaveDivider(mockRes, 10, "", "")

	// we expect total 2 deletion calls after all code got carried out
	mockDal.On("Delete", mock.Anything, mock.Anything).Return(nil).Twice()

	// for same type should return the same BatchSave
	jiraIssue1, err := divider.ForType(reflect.TypeOf(&MockJirIssueBsd{}))
	assert.Nil(t, err)

	jiraIssue2, err := divider.ForType(reflect.TypeOf(&MockJirIssueBsd{}))
	assert.Nil(t, err)
	assert.Equal(t, jiraIssue1, jiraIssue2)

	// for different types should return different BatchSaves
	jiraChangelog1, err := divider.ForType(reflect.TypeOf(&MockJiraChangelogBsd{}))
	assert.Nil(t, err)

	jiraChangelog2, err := divider.ForType(reflect.TypeOf(&MockJiraChangelogBsd{}))
	assert.Nil(t, err)
	assert.Equal(t, jiraChangelog1, jiraChangelog2)

	// assertion
	mockDal.AssertExpectations(t)
}
