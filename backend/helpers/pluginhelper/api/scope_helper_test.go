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
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/unithelper"
	mockcontext "github.com/apache/incubator-devlake/mocks/core/context"
	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type TestModel struct {
	ID   uint   `gorm:"primaryKey" validate:"required"`
	Name string `gorm:"primaryKey;type:BIGINT  NOT NULL"`
}

type TestFakeGitlabRepo struct {
	ConnectionId     uint64     `json:"connectionId" mapstructure:"connectionId" gorm:"primaryKey"`
	GitlabId         int        `json:"gitlabId" mapstructure:"gitlabId" gorm:"primaryKey"`
	CreatedDate      *time.Time `json:"createdDate" mapstructure:"-"`
	UpdatedDate      *time.Time `json:"updatedDate" mapstructure:"-"`
	common.NoPKModel `json:"-" mapstructure:"-"`
}

func (t TestFakeGitlabRepo) ScopeId() string {
	return fmt.Sprintf("%d", t.GitlabId)
}

func (t TestFakeGitlabRepo) ScopeName() string {
	return ""
}

func (t TestFakeGitlabRepo) ScopeFullName() string {
	return ""
}

func (t TestFakeGitlabRepo) TableName() string {
	return ""
}

func (t TestFakeGitlabRepo) ScopeParams() interface{} {
	return nil
}

type TestFakeGithubRepo struct {
	ConnectionId     uint64     `json:"connectionId" gorm:"primaryKey" mapstructure:"connectionId,omitempty"`
	GithubId         int        `json:"githubId" gorm:"primaryKey" mapstructure:"githubId"`
	Name             string     `json:"name" gorm:"type:varchar(255)" mapstructure:"name,omitempty"`
	HTMLUrl          string     `json:"HTMLUrl" gorm:"type:varchar(255)" mapstructure:"HTMLUrl,omitempty"`
	Description      string     `json:"description" mapstructure:"description,omitempty"`
	ScopeConfigId    uint64     `json:"scopeConfigId,omitempty" mapstructure:"scopeConfigId,omitempty"`
	OwnerId          int        `json:"ownerId" mapstructure:"ownerId,omitempty"`
	Language         string     `json:"language" gorm:"type:varchar(255)" mapstructure:"language,omitempty"`
	ParentGithubId   int        `json:"parentId" mapstructure:"parentGithubId,omitempty"`
	ParentHTMLUrl    string     `json:"parentHtmlUrl" mapstructure:"parentHtmlUrl,omitempty"`
	CloneUrl         string     `json:"cloneUrl" gorm:"type:varchar(255)" mapstructure:"cloneUrl,omitempty"`
	CreatedDate      *time.Time `json:"createdDate" mapstructure:"-"`
	UpdatedDate      *time.Time `json:"updatedDate" mapstructure:"-"`
	common.NoPKModel `json:"-" mapstructure:"-"`
}

func (r TestFakeGithubRepo) ScopeId() string {
	return fmt.Sprintf("%d", r.GithubId)
}

func (r TestFakeGithubRepo) ScopeName() string {
	return r.Name
}

func (r TestFakeGithubRepo) ScopeFullName() string {
	return r.Name
}

func (r TestFakeGithubRepo) ScopeParams() interface{} {
	return nil
}

func (TestFakeGithubRepo) TableName() string {
	return "_tool_github_repos"
}

type TestConnection struct {
	common.Model
	Name             string `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	Endpoint         string `mapstructure:"endpoint" env:"GITHUB_ENDPOINT" validate:"required"`
	Proxy            string `mapstructure:"proxy" env:"GITHUB_PROXY"`
	RateLimitPerHour int    `comment:"api request rate limit per hour"`
	Token            string `mapstructure:"token" env:"GITHUB_AUTH" validate:"required" encrypt:"yes"`
}

func (TestConnection) TableName() string {
	return "_tool_github_connections"
}

func TestVerifyScope(t *testing.T) {
	apiHelper := createMockScopeHelper[TestFakeGithubRepo]("GithubId")
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
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		err := apiHelper.verifyScope(&tc.model, validator.New())
		if (err != nil) != tc.wantErr {
			t.Errorf("unexpected error value - got: %v, want: %v", err, tc.wantErr)
		}

	}
}

type TestScopeConfig struct {
	common.Model         `mapstructure:"-"`
	Name                 string          `mapstructure:"name" json:"name" gorm:"type:varchar(255);index:idx_name_github,unique" validate:"required"`
	PrType               string          `mapstructure:"prType,omitempty" json:"prType" gorm:"type:varchar(255)"`
	PrComponent          string          `mapstructure:"prComponent,omitempty" json:"prComponent" gorm:"type:varchar(255)"`
	PrBodyClosePattern   string          `mapstructure:"prBodyClosePattern,omitempty" json:"prBodyClosePattern" gorm:"type:varchar(255)"`
	IssueSeverity        string          `mapstructure:"issueSeverity,omitempty" json:"issueSeverity" gorm:"type:varchar(255)"`
	IssuePriority        string          `mapstructure:"issuePriority,omitempty" json:"issuePriority" gorm:"type:varchar(255)"`
	IssueComponent       string          `mapstructure:"issueComponent,omitempty" json:"issueComponent" gorm:"type:varchar(255)"`
	IssueTypeBug         string          `mapstructure:"issueTypeBug,omitempty" json:"issueTypeBug" gorm:"type:varchar(255)"`
	IssueTypeIncident    string          `mapstructure:"issueTypeIncident,omitempty" json:"issueTypeIncident" gorm:"type:varchar(255)"`
	IssueTypeRequirement string          `mapstructure:"issueTypeRequirement,omitempty" json:"issueTypeRequirement" gorm:"type:varchar(255)"`
	DeploymentPattern    string          `mapstructure:"deploymentPattern,omitempty" json:"deploymentPattern" gorm:"type:varchar(255)"`
	ProductionPattern    string          `mapstructure:"productionPattern,omitempty" json:"productionPattern" gorm:"type:varchar(255)"`
	Refdiff              json.RawMessage `mapstructure:"refdiff,omitempty" json:"refdiff" swaggertype:"object" format:"json"`
}

func (TestScopeConfig) TableName() string {
	return "_tool_github_scope_configs"
}

func TestSetScopeFields(t *testing.T) {
	// create a struct
	type P struct {
		ConnectionId     uint64     `json:"connectionId" mapstructure:"connectionId" gorm:"primaryKey"`
		GitlabId         int        `json:"gitlabId" mapstructure:"gitlabId" gorm:"primaryKey"`
		CreatedAt        *time.Time `json:"createdAt" mapstructure:"-"`
		UpdatedAt        *time.Time `json:"updatedAt" mapstructure:"-"`
		common.NoPKModel `json:"-" mapstructure:"-"`
	}
	p := P{}
	apiHelper := createMockScopeHelper[TestFakeGitlabRepo]("GitlabId")

	// call setScopeFields to assign value
	connectionId := uint64(123)
	createdAt := time.Now()
	updatedAt := &createdAt
	apiHelper.setScopeFields(&p, connectionId, &createdAt, updatedAt)

	// verify fields
	if p.ConnectionId != connectionId {
		t.Errorf("ConnectionId not set correctly, expected: %v, got: %v", connectionId, p.ConnectionId)
	}

	if !p.CreatedAt.Equal(createdAt) {
		t.Errorf("CreatedAt not set correctly, expected: %v, got: %v", createdAt, p.CreatedAt)
	}

	if p.UpdatedAt == nil {
		t.Errorf("UpdatedDate not set correctly, expected: %v, got: %v", updatedAt, p.UpdatedAt)
	} else if !p.UpdatedAt.Equal(*updatedAt) {
		t.Errorf("UpdatedDate not set correctly, expected: %v, got: %v", updatedAt, p.UpdatedAt)
	}

	apiHelper.setScopeFields(&p, connectionId, &createdAt, nil)

	// verify fields
	if p.ConnectionId != connectionId {
		t.Errorf("ConnectionId not set correctly, expected: %v, got: %v", connectionId, p.ConnectionId)
	}

	if !p.CreatedAt.Equal(createdAt) {
		t.Errorf("CreatedDate not set correctly, expected: %v, got: %v", createdAt, p.CreatedAt)
	}

	if p.UpdatedAt != nil {
		t.Errorf("UpdatedDate not set correctly, expected: %v, got: %v", nil, p.UpdatedAt)
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
	result := returnPrimaryKeyValue(test)
	expected := "1-123"
	if result != expected {
		t.Errorf("returnPrimaryKeyValue returned %s, expected %s", result, expected)
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

	result2 := returnPrimaryKeyValue(test2)
	expected2 := ""
	if result2 != expected2 {
		t.Errorf("returnPrimaryKeyValue returned %s, expected %s", result2, expected2)
	}
}

func TestScopeApiHelper_Put(t *testing.T) {
	apiHelper := createMockScopeHelper[TestFakeGithubRepo]("GithubId")
	// create a mock input, scopes, and connection
	input := &plugin.ApiResourceInput{Params: map[string]string{"connectionId": "123"}, Body: map[string]interface{}{
		"data": []map[string]interface{}{
			{
				"HTMLUrl":       "string",
				"githubId":      1,
				"cloneUrl":      "string",
				"connectionId":  1,
				"createdAt":     "string",
				"createdDate":   "string",
				"description":   "string",
				"language":      "string",
				"name":          "string",
				"owner":         "string",
				"scopeConfigId": 0,
				"updatedAt":     "string",
				"updatedDate":   "string",
			},
			{
				"HTMLUrl":       "11",
				"githubId":      2,
				"cloneUrl":      "string",
				"connectionId":  1,
				"createdAt":     "string",
				"createdDate":   "string",
				"description":   "string",
				"language":      "string",
				"name":          "string",
				"owner":         "string",
				"scopeConfigId": 0,
				"updatedAt":     "string",
				"updatedDate":   "string",
			}}}}
	// test a successful call to Put
	_, err := apiHelper.Put(input)
	assert.NoError(t, err)
}

func createMockScopeHelper[Repo plugin.ToolLayerScope](scopeIdFieldName string) *ScopeApiHelper[TestConnection, Repo, TestScopeConfig] {
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
	mockDal.On("AllTables").Return(nil, nil)

	connHelper := NewConnectionHelper(mockRes, nil, "dummy_plugin")

	params := &ReflectionParameters{
		ScopeIdFieldName:  scopeIdFieldName,
		ScopeIdColumnName: "scope_id",
		RawScopeParamName: "ScopeId",
	}
	dbHelper := NewScopeDatabaseHelperImpl[TestConnection, Repo, TestScopeConfig](mockRes, connHelper, params)
	// create a mock ScopeApiHelper with a mock database connection
	return NewScopeHelper[TestConnection, Repo, TestScopeConfig](mockRes, validator.New(), connHelper, dbHelper, params, nil)
}
