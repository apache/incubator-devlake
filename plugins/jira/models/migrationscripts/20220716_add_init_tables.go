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
	"encoding/base64"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/jira/models/migrationscripts/archived"
)

type jiraConnectionNew struct {
	archived.RestConnection `mapstructure:",squash"`
	archived.BasicAuth      `mapstructure:",squash"`
}

func (jiraConnectionNew) TableName() string {
	return "_tool_jira_connections_new"
}

const JiraConnectionV011 = "_tool_jira_connections"

type jiraConnectionV011Old struct {
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

func (jiraConnectionV011Old) TableName() string {
	return "_tool_jira_connections_old"
}

type addInitTables struct{}

func (*addInitTables) Up(basicRes core.BasicRes) errors.Error {
	var err errors.Error
	if err = basicRes.GetDal().DropTables(
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
	); err != nil {
		return err
	}
	dal := basicRes.GetDal()

	// step1
	if err = migrationhelper.AutoMigrateTables(basicRes, &jiraConnectionNew{}); err != nil {
		return err
	}
	//nolint:errcheck
	defer dal.DropTables(&jiraConnectionNew{})

	// step2
	if err = dal.RenameTable("_tool_jira_connections", "_tool_jira_connections_old"); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = dal.RenameTable("_tool_jira_connections_old", "_tool_jira_connections")
		} else {
			err = dal.DropTables(&jiraConnectionV011Old{})
		}
	}()

	// step3
	var jiraConns []jiraConnectionV011Old
	err = dal.All(&jiraConns)
	if err != nil {
		return err
	}

	for _, v := range jiraConns {
		conn := &jiraConnectionNew{}
		conn.ID = v.ID
		conn.Name = v.Name
		conn.Endpoint = v.Endpoint
		conn.Proxy = v.Proxy
		conn.RateLimitPerHour = v.RateLimit

		c := config.GetConfig()
		encKey := c.GetString(core.EncodeKeyEnvStr)
		if encKey == "" {
			return errors.BadInput.New("jira v0.11 invalid encKey")
		}
		var auth string
		if auth, err = core.Decrypt(encKey, v.BasicAuthEncoded); err != nil {
			return errors.Convert(err)
		}
		//var pk []byte
		pk, err1 := base64.StdEncoding.DecodeString(auth)
		if err1 != nil {
			return errors.Convert(err1)
		}
		originInfo := strings.Split(string(pk), ":")
		if len(originInfo) == 2 {
			conn.Username = originInfo[0]
			conn.Password, err = core.Encrypt(encKey, originInfo[1])
			if err != nil {
				return err
			}
			// create
			err = dal.Create(&conn)
			if err != nil {
				return err
			}
		}
	}

	// step4
	err = dal.RenameTable("_tool_jira_connections_new", "_tool_jira_connections")
	if err != nil {
		return errors.Convert(err)
	}

	return errors.Convert(migrationhelper.AutoMigrateTables(
		basicRes,
		&archived.JiraAccount{},
		&archived.JiraBoardIssue{},
		&archived.JiraBoard{},
		&archived.JiraIssueChangelogItems{},
		&archived.JiraIssueChangelogs{},
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
	))
}

func (*addInitTables) Version() uint64 {
	return 20220716201138
}

func (*addInitTables) Name() string {
	return "Jira init schemas"
}
