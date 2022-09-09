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

type JiraConnectionNew struct {
	archived.RestConnection `mapstructure:",squash"`
	archived.BasicAuth      `mapstructure:",squash"`
}

func (JiraConnectionNew) TableName() string {
	return "_tool_jira_connections_new"
}

const JiraConnectionV011 = "_tool_jira_connections"

type JiraConnectionV011Old struct {
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

func (JiraConnectionV011Old) TableName() string {
	return "_tool_jira_connections_old"
}

type addInitTables struct{}

func (*addInitTables) Up(ctx context.Context, db *gorm.DB) (err error) {

	err = db.Migrator().DropTable(
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

	// In order to avoid postgres reporting "cached plan must not change result type",
	// we have to avoid loading the _tool_jira_connections table before our migration. The trick is to rename it before loading.
	// so need to use JiraConnectionOld, JiraConnectionNew and JiraConnection to operate.
	// step1: create _tool_jira_connections_new table
	// step2: rename _tool_jira_connections to _tool_jira_connections_old
	// step3: select data from _tool_jira_connections_old and insert date to _tool_jira_connections_new
	// step4: rename _tool_jira_connections_new to _tool_jira_connections

	// step1
	err = db.Migrator().CreateTable(&JiraConnectionNew{})
	if err != nil {
		return err
	}
	defer db.Migrator().DropTable(&JiraConnectionNew{})

	// step2
	err = db.Migrator().RenameTable("_tool_jira_connections", "_tool_jira_connections_old")
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = db.Migrator().RenameTable("_tool_jira_connections_old", "_tool_jira_connections")
		} else {
			err = db.Migrator().DropTable(&JiraConnectionV011Old{})
		}
	}()

	// step3
	var result *gorm.DB
	var jiraConns []JiraConnectionV011Old
	result = db.Find(&jiraConns)
	if result.Error != nil {
		return result.Error
	}

	for _, v := range jiraConns {
		conn := &JiraConnectionNew{}
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
			tx := db.Create(&conn)
			if tx.Error != nil {
				return tx.Error
			}
		}
	}

	// step4
	err = db.Migrator().RenameTable("_tool_jira_connections_new", "_tool_jira_connections")
	if err != nil {
		return err
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

func (*addInitTables) Version() uint64 {
	return 20220716201138
}

func (*addInitTables) Name() string {
	return "Jira init schemas"
}
