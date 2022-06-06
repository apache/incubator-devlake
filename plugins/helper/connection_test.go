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
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type MockConnection struct {
	RestConnection             `mapstructure:",squash"`
	BasicAuth                  `mapstructure:",squash"`
	StoryPointField            string `gorm:"type:varchar(50);" json:"storyPointField"`
	RemotelinkCommitShaPattern string `gorm:"type:varchar(255);comment='golang regexp, the first group will be recognized as commit sha, ref https://github.com/google/re2/wiki/Syntax'" json:"remotelinkCommitShaPattern"`
}

func (MockConnection) TableName() string {
	return "_tool_jira_connections"
}

func TestMergeFieldsToConnection(t *testing.T) {
	v := &MockConnection{
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
		RemotelinkCommitShaPattern: "8",
	}
	data := make(map[string]interface{})
	data["Endpoint"] = "2-2"
	data["Username"] = "4-4"
	data["Password"] = "5-5"

	err := mergeFieldsToConnection(v, data)
	assert.Nil(t, err)

	assert.Equal(t, "4-4", v.Username)
	assert.Equal(t, "2-2", v.Endpoint)
	assert.Equal(t, "5-5", v.Password)
}

func TestDecryptAndEncrypt(t *testing.T) {
	v := &MockConnection{
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
		RemotelinkCommitShaPattern: "8",
	}
	err := EncryptConnection(v)
	assert.Nil(t, err)

	assert.NotEqual(t, "5", v.Password)
	err = DecryptConnection(v)
	assert.Nil(t, err)

	assert.Equal(t, "5", v.Password)

}

func TestDecryptConnection(t *testing.T) {
	v := &MockConnection{
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
		RemotelinkCommitShaPattern: "8",
	}
	err := EncryptConnection(v)
	assert.Nil(t, err)

	encryptedPwd := v.Password
	err = DecryptConnection(v)
	assert.Nil(t, err)
	assert.NotEqual(t, encryptedPwd, v.Password)
	assert.Equal(t, "5", v.Password)
}

func TestFirstFieldNameWithTag(t *testing.T) {
	v := &MockConnection{
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
		StoryPointField:            "7",
		RemotelinkCommitShaPattern: "8",
	}
	dataVal := reflect.ValueOf(v)
	dataType := reflect.Indirect(dataVal).Type()
	fieldName := firstFieldNameWithTag(dataType, "encrypt")
	assert.Equal(t, "Password", fieldName)
}

//func TestListConnections(t *testing.T) {
//	jiraConnections := make([]*MockConnection, 0)
//	cfg := config.GetConfig()
//	dbUrl := cfg.GetString("DB_URL")
//	u, err := url.Parse(dbUrl)
//	dbUrl = fmt.Sprintf("%s@tcp(%s)%s?%s", u.User.String(), u.Host, u.Path, u.RawQuery)
//	dbConfig := &gorm.Config{
//		Logger: gormLogger.New(
//			log.Default(),
//			gormLogger.Config{
//				SlowThreshold:             time.Second,      // Slow SQL threshold
//				LogLevel:                  gormLogger.Error, // Log level
//				IgnoreRecordNotFoundError: true,             // Ignore ErrRecordNotFound error for logger
//				Colorful:                  true,             // Disable color
//			},
//		),
//		// most of our operation are in batch, this can improve performance
//		PrepareStmt: true,
//	}
//	db, err := gorm.Open(mysql.Open(dbUrl), dbConfig)
//
//	err = ListConnections(&jiraConnections, db)
//
//	assert.Nil(t, err)
//}
