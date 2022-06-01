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
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/plugins/core"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestConnection struct {
	RestConnection             `mapstructure:",squash"`
	BasicAuth                  `mapstructure:",squash"`
	EpicKeyField               string `gorm:"type:varchar(50);" json:"epicKeyField"`
	StoryPointField            string `gorm:"type:varchar(50);" json:"storyPointField"`
	RemotelinkCommitShaPattern string `gorm:"type:varchar(255);comment='golang regexp, the first group will be recognized as commit sha, ref https://github.com/google/re2/wiki/Syntax'" json:"remotelinkCommitShaPattern"`
}

func TestMergeFieldsToConnection(t *testing.T) {
	v := &TestConnection{
		RestConnection: RestConnection{
			BaseConnection: BaseConnection{
				Name: "1",
			},
			Endpoint:  "2",
			Proxy:     "3",
			RateLimit: 0,
		},
		BasicAuth: BasicAuth{
			Username: "4",
			Password: "5",
		},
		EpicKeyField:               "6",
		StoryPointField:            "7",
		RemotelinkCommitShaPattern: "8",
	}
	data := make(map[string]interface{})
	data["Endpoint"] = "2-2"
	data["Username"] = "4-4"
	data["Password"] = "5-5"

	err := mergeFieldsToConnection(v, data)
	if err != nil {
		return
	}

	assert.Equal(t, "4-4", v.Username)
	assert.Equal(t, "2-2", v.Endpoint)
	assert.Equal(t, "5-5", v.Password)
}

func TestDecryptAndEncrypt(t *testing.T) {
	v := &TestConnection{
		RestConnection: RestConnection{
			BaseConnection: BaseConnection{
				Name: "1",
			},
			Endpoint:  "2",
			Proxy:     "3",
			RateLimit: 0,
		},
		BasicAuth: BasicAuth{
			Username: "4",
			Password: "5",
		},
		EpicKeyField:               "6",
		StoryPointField:            "7",
		RemotelinkCommitShaPattern: "8",
	}
	dataVal := reflect.ValueOf(v)
	encKey := "test"
	err := encryptField(dataVal, "Password", encKey)
	if err != nil {
		return
	}
	assert.NotEqual(t, "5", v.Password)
	err = decryptField(dataVal, "Password", encKey)
	if err != nil {
		return
	}

	assert.Equal(t, "5", v.Password)

}

func TestDecryptConnection(t *testing.T) {
	v := &TestConnection{
		RestConnection: RestConnection{
			BaseConnection: BaseConnection{
				Name: "1",
			},
			Endpoint:  "2",
			Proxy:     "3",
			RateLimit: 0,
		},
		BasicAuth: BasicAuth{
			Username: "4",
			Password: "5",
		},
		EpicKeyField:               "6",
		StoryPointField:            "7",
		RemotelinkCommitShaPattern: "8",
	}
	encKey, err := getEncKey()
	if err != nil {
		return
	}
	dataVal := reflect.ValueOf(v)
	err = encryptField(dataVal, "Password", encKey)
	if err != nil {
		return
	}
	encryptedPwd := v.Password
	err = DecryptConnection(v, "Password")
	if err != nil {
		return
	}
	assert.NotEqual(t, encryptedPwd, v.Password)
	assert.Equal(t, "5", v.Password)
}

func TestGetEncKey(t *testing.T) {
	// encryptField
	v := config.GetConfig()
	encKey := v.GetString(core.EncodeKeyEnvStr)
	str, err := getEncKey()
	if err != nil {
		return
	}
	if len(encKey) > 0 {
		assert.Equal(t, encKey, str)
	} else {
		assert.NotEqual(t, 0, len(str))
	}

}

func TestFirstFieldNameWithTag(t *testing.T) {
	v := &TestConnection{
		RestConnection: RestConnection{
			BaseConnection: BaseConnection{
				Name: "1",
			},
			Endpoint:  "2",
			Proxy:     "3",
			RateLimit: 0,
		},
		BasicAuth: BasicAuth{
			Username: "4",
			Password: "5",
		},
		EpicKeyField:               "6",
		StoryPointField:            "7",
		RemotelinkCommitShaPattern: "8",
	}
	dataVal := reflect.ValueOf(v)
	dataType := reflect.Indirect(dataVal).Type()
	fieldName := firstFieldNameWithTag(dataType, "encryptField")
	assert.Equal(t, "Password", fieldName)
}
