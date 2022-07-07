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

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/jira/models/migrationscripts/archived"
	"gorm.io/gorm"
)

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

	if m.HasTable(&archived.JiraConnectionV11{}) {
		var jiraConns []archived.JiraConnectionV11
		result = db.Find(&jiraConns)

		if result.Error == nil {
			err := db.Migrator().DropTable(&archived.JiraConnectionV11{})
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
				conn.RateLimit = v.RateLimit

				c := config.GetConfig()
				encKey := c.GetString("ENCODE_KEY")
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
					conn.Password = originInfo[1]
					// create
					db.Create(&conn)
				}
			}
		} else if m.HasTable(&archived.JiraConnectionV10{}) {
			var jiraConns []archived.JiraConnectionV10
			result = db.Find(&jiraConns)

			if result.Error == nil {
				err := db.Migrator().DropTable(&archived.JiraConnectionV10{})
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
					conn.RateLimit = v.RateLimit

					c := config.GetConfig()
					encKey := c.GetString("ENCODE_KEY")
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
						conn.Password = originInfo[1]
						// create
						db.Create(&conn)
					}
				}
			}
		} else {
			return result.Error
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
	)
}

func (*InitSchemas) Version() uint64 {
	return 20220707201138
}

func (*InitSchemas) Name() string {
	return "Jira init schemas"
}
