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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/unithelper"
	mockcontext "github.com/apache/incubator-devlake/mocks/core/context"
	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"reflect"
	"testing"
	"time"
)

type TestModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"primaryKey;type:BIGINT  NOT NULL"`
}

type GithubRepo struct {
	ConnectionId         uint64     `json:"connectionId" gorm:"primaryKey" mapstructure:"connectionId,omitempty"`
	GithubId             int        `json:"githubId" gorm:"primaryKey" mapstructure:"githubId"`
	Name                 string     `json:"name" gorm:"type:varchar(255)" mapstructure:"name,omitempty"`
	HTMLUrl              string     `json:"HTMLUrl" gorm:"type:varchar(255)" mapstructure:"HTMLUrl,omitempty"`
	Description          string     `json:"description" mapstructure:"description,omitempty"`
	TransformationRuleId uint64     `json:"transformationRuleId,omitempty" mapstructure:"transformationRuleId,omitempty"`
	OwnerId              int        `json:"ownerId" mapstructure:"ownerId,omitempty"`
	Language             string     `json:"language" gorm:"type:varchar(255)" mapstructure:"language,omitempty"`
	ParentGithubId       int        `json:"parentId" mapstructure:"parentGithubId,omitempty"`
	ParentHTMLUrl        string     `json:"parentHtmlUrl" mapstructure:"parentHtmlUrl,omitempty"`
	CloneUrl             string     `json:"cloneUrl" gorm:"type:varchar(255)" mapstructure:"cloneUrl,omitempty"`
	CreatedDate          *time.Time `json:"createdDate" mapstructure:"-"`
	UpdatedDate          *time.Time `json:"updatedDate" mapstructure:"-"`
	common.NoPKModel     `json:"-" mapstructure:"-"`
}

func (GithubRepo) TableName() string {
	return "_tool_github_repos"
}

type GithubConnection struct {
	common.Model
	Name             string `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	Endpoint         string `mapstructure:"endpoint" env:"GITHUB_ENDPOINT" validate:"required"`
	Proxy            string `mapstructure:"proxy" env:"GITHUB_PROXY"`
	RateLimitPerHour int    `comment:"api request rate limit per hour"`
	Token            string `mapstructure:"token" env:"GITHUB_AUTH" validate:"required" encrypt:"yes"`
}

func (GithubConnection) TableName() string {
	return "_tool_github_connections"
}

type req struct {
	Data []*GithubRepo `json:"data"`
}

func TestCheckPrimaryKeyValue(t *testing.T) {
	testCases := []struct {
		name    string
		model   TestModel
		wantErr bool
	}{
		{
			name: "valid case",
			model: TestModel{
				ID:   1,
				Name: "test name",
			},
			wantErr: false,
		},
		{
			name: "zero value",
			model: TestModel{
				ID:   0,
				Name: "test name",
			},
			wantErr: true,
		},
		{
			name: "nil value",
			model: TestModel{
				ID:   1,
				Name: "",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		err := VerifyPrimaryKeyValue(tc.model)
		if (err != nil) != tc.wantErr {
			t.Errorf("unexpected error value - got: %v, want: %v", err, tc.wantErr)
		}

	}
}

func TestSetGitlabProjectFields(t *testing.T) {
	// create a struct
	var p struct {
		ConnectionId uint64 `json:"connectionId" mapstructure:"connectionId" gorm:"primaryKey"`
		GitlabId     int    `json:"gitlabId" mapstructure:"gitlabId" gorm:"primaryKey"`

		CreatedDate      *time.Time `json:"createdDate" mapstructure:"-"`
		UpdatedDate      *time.Time `json:"updatedDate" mapstructure:"-"`
		common.NoPKModel `json:"-" mapstructure:"-"`
	}

	// call SetScopeFields to assign value
	connectionId := uint64(123)
	createdDate := time.Now()
	updatedDate := &createdDate
	SetScopeFields(&p, connectionId, &createdDate, updatedDate)

	// verify fields
	if p.ConnectionId != connectionId {
		t.Errorf("ConnectionId not set correctly, expected: %v, got: %v", connectionId, p.ConnectionId)
	}

	if !p.CreatedDate.Equal(createdDate) {
		t.Errorf("CreatedDate not set correctly, expected: %v, got: %v", createdDate, p.CreatedDate)
	}

	if p.UpdatedDate == nil {
		t.Errorf("UpdatedDate not set correctly, expected: %v, got: %v", updatedDate, p.UpdatedDate)
	} else if !p.UpdatedDate.Equal(*updatedDate) {
		t.Errorf("UpdatedDate not set correctly, expected: %v, got: %v", updatedDate, p.UpdatedDate)
	}

	SetScopeFields(&p, connectionId, &createdDate, nil)

	// verify fields
	if p.ConnectionId != connectionId {
		t.Errorf("ConnectionId not set correctly, expected: %v, got: %v", connectionId, p.ConnectionId)
	}

	if !p.CreatedDate.Equal(createdDate) {
		t.Errorf("CreatedDate not set correctly, expected: %v, got: %v", createdDate, p.CreatedDate)
	}

	if p.UpdatedDate != nil {
		t.Errorf("UpdatedDate not set correctly, expected: %v, got: %v", nil, p.UpdatedDate)
	}
}

func TestReturnPrimaryKeyValue(t *testing.T) {
	// Define a test struct with the primaryKey tag on one of its fields.
	type TestStruct struct {
		ConnectionId int    `json:"connectionId" gorm:"primaryKey"`
		Id           int    `json:"id" gorm:"primaryKey"`
		Name         string `json:"name"`
		CreatedAt    time.Time
		UpdatedAt    time.Time
		DeletedAt    gorm.DeletedAt `gorm:"index"`
	}

	// Create an instance of the test struct.
	test := TestStruct{
		ConnectionId: 1,
		Id:           123,
		Name:         "Test",
		CreatedAt:    time.Now(),
	}

	// Call the function and check if it returns the correct primary key value.
	result := ReturnPrimaryKeyValue(test)
	expected := "1-123"
	if result != expected {
		t.Errorf("ReturnPrimaryKeyValue returned %s, expected %s", result, expected)
	}

	// Test with a different struct that has no field with primaryKey tag.
	type TestStruct2 struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	test2 := TestStruct2{
		Id:        456,
		Name:      "Test 2",
		CreatedAt: time.Now(),
	}

	result2 := ReturnPrimaryKeyValue(test2)
	expected2 := ""
	if result2 != expected2 {
		t.Errorf("ReturnPrimaryKeyValue returned %s, expected %s", result2, expected2)
	}
}

func TestScopeApiHelper_Put(t *testing.T) {
	mockDal := new(mockdal.Dal)
	mockLogger := unithelper.DummyLogger()
	mockRes := new(mockcontext.BasicRes)

	mockRes.On("GetDal").Return(mockDal)
	mockRes.On("GetConfig", mock.Anything).Return("")
	mockRes.On("GetLogger").Return(mockLogger)

	// we expect total 2 deletion calls after all code got carried out
	mockDal.On("Delete", mock.Anything, mock.Anything).Return(nil).Twice()
	mockDal.On("GetPrimaryKeyFields", mock.Anything).Return(
		[]reflect.StructField{
			{Name: "ID", Type: reflect.TypeOf("")},
		},
	)
	mockDal.On("CreateOrUpdate", mock.Anything, mock.Anything).Return(nil)
	mockDal.On("First", mock.Anything, mock.Anything).Return(nil)
	mockDal.On("All", mock.Anything, mock.Anything).Return(nil)

	connHelper := NewConnectionHelper(mockRes, nil)

	// create a mock input, scopes, and connection
	input := &plugin.ApiResourceInput{Params: map[string]string{"connectionId": "123"}}
	scopes := []*GithubRepo{
		{GithubId: 1, Name: "scope1"},
		{GithubId: 2, Name: "scope2"},
	}
	connection := &GithubConnection{}
	connection.ID = 3
	connection.Name = "test"

	// create a mock ScopeApiHelper with a mock database connection
	apiHelper := &ScopeApiHelper{db: mockDal, connHelper: connHelper}
	apiReq := req{Data: scopes}
	// test a successful call to Put
	err := apiHelper.Put(input, &apiReq, connection)
	assert.NoError(t, err)

	// test a call to Put with a duplicate primary key value
	duplicateScopes := []*GithubRepo{
		{GithubId: 1, Name: "scope1"},
		{GithubId: 1, Name: "scope2"},
	}
	apiReq.Data = duplicateScopes
	err = apiHelper.Put(input, &apiReq, connection)
	assert.Error(t, err)

	// test a call to Put with an invalid primary key value
	invalidScopes := []*GithubRepo{
		{GithubId: 0, Name: "scope1"},
	}
	apiReq.Data = invalidScopes

	err = apiHelper.Put(input, &apiReq, connection)
	assert.Error(t, err)
}
