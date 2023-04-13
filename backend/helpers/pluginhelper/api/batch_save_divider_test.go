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

package api

import (
	"reflect"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/helpers/unithelper"
	mockcontext "github.com/apache/incubator-devlake/mocks/core/context"
	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"

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
	mockDal := new(mockdal.Dal)

	mockLogger := unithelper.DummyLogger()
	mockRes := new(mockcontext.BasicRes)

	mockRes.On("GetDal").Return(mockDal)
	mockRes.On("GetLogger").Return(mockLogger)

	// we expect total 2 deletion calls after all code got carried out
	mockDal.On("Delete", mock.Anything, mock.Anything).Return(nil).Twice()
	mockDal.On("GetPrimaryKeyFields", mock.Anything).Return(
		[]reflect.StructField{
			{Name: "ID", Type: reflect.TypeOf("")},
		},
	)

	divider := NewBatchSaveDivider(mockRes, 10, "a", "b")

	// for same type should return the same BatchSave
	jiraIssue1, err := divider.ForType(reflect.TypeOf(&MockJirIssueBsd{}))
	assert.Nil(t, err)

	jiraIssue2, err := divider.ForType(reflect.TypeOf(&MockJirIssueBsd{}))
	assert.Nil(t, err)
	assert.Equal(t, jiraIssue1, jiraIssue2)

	// for different types should return different BatchSaves
	jiraChangelog1, err := divider.ForType(reflect.TypeOf(&MockJiraChangelogBsd{}))
	assert.Nil(t, err)

	jiraChangelog, err := divider.ForType(reflect.TypeOf(&MockJiraChangelogBsd{}))
	assert.Nil(t, err)
	assert.Equal(t, jiraChangelog1, jiraChangelog)

	// assertion
	mockDal.AssertExpectations(t)
}
