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
	"fmt"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/jira/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type JiraConnectionV011 struct {
	ID                         uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt                  time.Time `json:"createdAt"`
	UpdatedAt                  time.Time `json:"updatedAt"`
	Name                       string    `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	Endpoint                   string    `json:"endpoint" validate:"required"`
	BasicAuthEncoded           string    `json:"basicAuthEncoded" validate:"required"`
	EpicKeyField               string    `gorm:"type:varchar(50);" json:"epicKeyField"`
	StoryPointField            string    `gorm:"type:varchar(50);" json:"storyPointField"`
	RemotelinkCommitShaPattern string    `gorm:"type:varchar(255);comment='golang regexp, the first group will be recognized as commit sha, ref https://github.com/google/re2/wiki/Syntax'" json:"remotelinkCommitShaPattern"`
	Proxy                      string    `json:"proxy"`
	RateLimit                  int       `comment:"api request rate limt per hour" json:"rateLimit"`
}

func (JiraConnectionV011) TableName() string {
	return "_tool_jira_connections"
}

type InitSchemas struct{}

func (*InitSchemas) Up(ctx context.Context, db *gorm.DB) error {

	err := db.Migrator().DropTable(
		// history table
		"_raw_jira_api_users",
		"_raw_jira_api_boards",
		"_raw_jira_api_changelogs",
		"_raw_jira_api_issues",
		"_raw_jira_api_projects",
		"_raw_jira_api_remotelinks",
		"_raw_jira_api_sprints",
		"_raw_jira_api_status",
		"_raw_jira_api_worklogs",
		"_tool_jira_accounts",
		"_tool_jira_issue_type_mappings",
		"_tool_jira_issue_status_mappings",
		"_tool_jira_changelogs",
		"_tool_jira_changelog_items",
		&archived.JiraProject{},
		&archived.JiraIssue{},
		&archived.JiraBoard{},
		&archived.JiraBoardIssue{},
		&archived.JiraRemotelink{},
		&archived.JiraIssueCommit{},
		&archived.JiraSprint{},
		&archived.JiraBoardSprint{},
		&archived.JiraSprintIssue{},
		&archived.JiraWorklog{},
	)
	if err != nil {
		return err
	}

	// get connection history data
	var result *gorm.DB
	m := db.Migrator()

	if m.HasTable(&JiraConnectionV011{}) {
		var jiraConns []JiraConnectionV011
		result = db.Find(&jiraConns)

		if result.Error == nil {
			err := db.Migrator().DropTable(&JiraConnectionV011{})
			if err != nil {
				return err
			}
			err = db.Migrator().AutoMigrate(&archived.JiraConnection{})
			if err != nil {
				return err
			}

			for _, v := range jiraConns {
				conn := &archived.JiraConnection{}
				conn.ID = v.ID
				conn.Name = v.Name
				conn.Endpoint = v.Endpoint
				conn.Proxy = v.Proxy
				conn.RateLimitPerHour = v.RateLimit

				c := config.GetConfig()
				encKey := c.GetString(core.EncodeKeyEnvStr)
				if encKey == "" {
					return fmt.Errorf("jira v0.11 invalid encKey")
				}
				//base64.StdEncoding.DecodeString() v.BasicAuthEncoded
				auth, err := core.Decrypt(encKey, v.BasicAuthEncoded)
				if err != nil {
					return err
				}
				pk, err := base64.StdEncoding.DecodeString(auth)
				if err != nil {
					return err
				}
				originInfo := strings.Split(string(pk), ":")
				if len(originInfo) == 2 {
					conn.Username = originInfo[0]
					conn.Password, err = core.Encrypt(encKey, originInfo[1])
					if err != nil {
						return err
					}
					// create
					db.Create(&conn)
				}
			}
		}
	} else {
		c := config.GetConfig()
		encKey := c.GetString("ENCODE_KEY")
		if encKey == "" {
			return fmt.Errorf("jira invalid encKey")
		}
		err := db.Migrator().AutoMigrate(&archived.JiraConnection{})
		if err != nil {
			return err
		}
	}

	return db.Migrator().AutoMigrate(
		&archived.JiraAccount{},
		&archived.JiraBoardIssue{},
		&archived.JiraBoard{},
		&archived.JiraIssueChangelogItems{},
		&archived.JiraIssueChangelogs{},
		//&archived.JiraConnection{},
		&archived.JiraIssueCommit{},
		&archived.JiraIssueLabel{},
		&archived.JiraIssue{},
		&archived.JiraProject{},
		&archived.JiraRemotelink{},
		&archived.JiraSprint{},
		&archived.JiraBoardSprint{},
		&archived.JiraSprintIssue{},
		&archived.JiraStatus{},
		&archived.JiraWorklog{},
		&archived.JiraIssueType{},
	)
}

func (*InitSchemas) Version() uint64 {
	return 20220716201138
}

func (*InitSchemas) Name() string {
	return "Jira init schemas"
}
