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
	"context"
	"encoding/base64"
	"strings"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"gorm.io/gorm"
)

type JiraConnection20220601 struct {
	helper.RestConnection
	helper.BasicAuth
	EpicKeyField               string `gorm:"type:varchar(50);" json:"epicKeyField"`
	StoryPointField            string `gorm:"type:varchar(50);" json:"storyPointField"`
	RemotelinkCommitShaPattern string `gorm:"type:varchar(255);comment='golang regexp, the first group will be recognized as commit sha, ref https://github.com/google/re2/wiki/Syntax'" json:"remotelinkCommitShaPattern"`
}

func (JiraConnection20220601) TableName() string {
	return "_tool_jira_connections"
}

type UpdateSchemas20220601 struct {
	config core.ConfigGetter
	logger core.Logger
}

func (u *UpdateSchemas20220601) SetConfigGetter(getter core.ConfigGetter) {
	u.config = getter
}

func (u *UpdateSchemas20220601) SetLogger(logger core.Logger) {
	u.logger = logger
}

func (u *UpdateSchemas20220601) Up(ctx context.Context, db *gorm.DB) error {
	var err error
	if !db.Migrator().HasColumn(&JiraConnection20220505{}, "password") {
		err = db.Migrator().AddColumn(&JiraConnection20220601{}, "password")
		if err != nil {
			return err
		}
	}

	if !db.Migrator().HasColumn(&JiraConnection20220505{}, "username") {
		err = db.Migrator().AddColumn(&JiraConnection20220601{}, "username")
		if err != nil {
			return err
		}
	}

	if db.Migrator().HasColumn(&JiraConnection20220505{}, "basic_auth_encoded") {
		connections := make([]*JiraConnection20220505, 0)
		err = db.Find(&connections).Error
		if err != nil {
			return err
		}
		encKey := u.config.GetString(core.EncodeKeyEnvStr)
		for _, connection := range connections {
			basicAuthEncoded, err := core.Decrypt(encKey, connection.BasicAuthEncoded)
			if err != nil {
				return err
			}
			basicAuth, err := base64.StdEncoding.DecodeString(basicAuthEncoded)
			if err != nil {
				return err
			}
			strList := strings.Split(string(basicAuth), ":")
			if len(strList) > 1 {
				encPass, err := core.Encrypt(encKey, strList[1])
				if err != nil {
					return err
				}
				newConnection := JiraConnection20220601{
					RestConnection: helper.RestConnection{
						BaseConnection: helper.BaseConnection{
							Name:  connection.Name,
							Model: connection.Model,
						},
						Endpoint:  connection.Endpoint,
						Proxy:     connection.Proxy,
						RateLimit: connection.RateLimit,
					},
					BasicAuth: helper.BasicAuth{
						Username: strList[0],
						Password: encPass,
					},
					EpicKeyField:               connection.EpicKeyField,
					StoryPointField:            connection.StoryPointField,
					RemotelinkCommitShaPattern: connection.RemotelinkCommitShaPattern,
				}
				err = db.Save(newConnection).Error
				if err != nil {
					return err
				}
			}
		}
		err = db.Migrator().DropColumn(&JiraConnection20220505{}, "basic_auth_encoded")
		if err != nil {
			return err
		}
	}

	return nil
}

func (*UpdateSchemas20220601) Version() uint64 {
	return 20220601154646
}

func (*UpdateSchemas20220601) Name() string {
	return "change basic_auth to username/password"
}
