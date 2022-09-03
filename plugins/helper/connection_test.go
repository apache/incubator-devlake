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
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockAuth struct {
	Username string
	Password string `encrypt:"yes"`
}

type MockConnection struct {
	MockAuth
	Name      string `mapstructure:"name"`
	BasicAuth string `encrypt:"true"`
	BearToken struct {
		AccessToken string `encrypt:"true"`
	}
	MockAuth2 *MockAuth
	Age       int
	Since     *time.Time
}

/*
func TestMergeFieldsToConnection(t *testing.T) {
	v := &MockConnection{
		Name: "1",
		BearToken: struct {
			AccessToken string "encrypt:\"true\""
		}{
			AccessToken: "2",
		},
		MockAuth: &MockAuth{
			Username: "3",
			Password: "4",
		},
		Age: 5,
	}
	data := make(map[string]interface{})
	data["name"] = "1a"
	data["BasicAuth"] = map[string]interface{}{
		"AccessToken": "2a",
	}
	data["Username"] = "3a"

	err := mergeFieldsToConnection(v, data)
	assert.Nil(t, err)

	assert.Equal(t, "1a", v.Name)
	assert.Equal(t, "2a", v.BearToken.AccessToken)
	assert.Equal(t, "3a", v.Username)
	assert.Equal(t, "4", v.Password)
	assert.Equal(t, 5, v.Age)
}
*/

func TestUpdateEncryptFields(t *testing.T) {
	sinc := time.Now()
	v := &MockConnection{
		MockAuth: MockAuth{
			Username: "1",
			Password: "2",
		},
		Name: "3",
		BearToken: struct {
			AccessToken string `encrypt:"true"`
		}{
			AccessToken: "4",
		},
		MockAuth2: &MockAuth{
			Username: "5",
			Password: "6",
		},
		Age:   7,
		Since: &sinc,
	}
	err := UpdateEncryptFields(v, func(in string) (string, errors.Error) {
		return fmt.Sprintf("%s-asdf", in), nil
	})
	assert.Nil(t, err)
	assert.Equal(t, "1", v.Username)
	assert.Equal(t, "2-asdf", v.Password)
	assert.Equal(t, "3", v.Name)
	assert.Equal(t, "4-asdf", v.BearToken.AccessToken)
	assert.Equal(t, "5", v.MockAuth2.Username)
	assert.Equal(t, "6-asdf", v.MockAuth2.Password)
	assert.Equal(t, 7, v.Age)
}
