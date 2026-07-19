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

package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/codequality"
	coremigrations "github.com/apache/incubator-devlake/core/models/migrationscripts"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	implcontext "github.com/apache/incubator-devlake/impls/context"
	"github.com/apache/incubator-devlake/plugins/sonarqube/impl"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
	sonarqubemigrations "github.com/apache/incubator-devlake/plugins/sonarqube/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/sonarqube/tasks"
	"github.com/stretchr/testify/require"
)

type sonarqubeIssueCodeBlockBeforeText struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	Id           string `gorm:"primaryKey"`
	IssueKey     string `gorm:"index"`
	Component    string `gorm:"index;type:varchar(500)"`
	StartLine    int
	EndLine      int
	StartOffset  int
	EndOffset    int
	Msg          string
	common.NoPKModel
}

func (sonarqubeIssueCodeBlockBeforeText) TableName() string {
	return "_tool_sonarqube_issue_code_blocks"
}

type cqIssueCodeBlockBeforeText struct {
	domainlayer.DomainEntity
	IssueKey    string `json:"key" gorm:"index"`
	Component   string `gorm:"index"`
	StartLine   int
	EndLine     int
	StartOffset int
	EndOffset   int
	Msg         string
}

func (cqIssueCodeBlockBeforeText) TableName() string {
	return "cq_issue_code_blocks"
}

func TestSonarqubeIssueCodeBlockLongComponent(t *testing.T) {
	var sonarqube impl.Sonarqube
	dataflowTester := e2ehelper.NewDataFlowTester(t, "sonarqube", sonarqube)
	dataflowTester.FlushTabler(&models.SonarqubeIssue{})
	require.NoError(t, dataflowTester.Db.Migrator().DropTable(
		&sonarqubeIssueCodeBlockBeforeText{},
		&cqIssueCodeBlockBeforeText{},
	))
	require.NoError(t, dataflowTester.Db.AutoMigrate(
		&sonarqubeIssueCodeBlockBeforeText{},
		&cqIssueCodeBlockBeforeText{},
	))

	existingComponent := "existing:component"
	require.NoError(t, dataflowTester.Db.Create(&sonarqubeIssueCodeBlockBeforeText{
		ConnectionId: 1,
		Id:           "existing-tool-block",
		IssueKey:     "existing-issue",
		Component:    existingComponent,
	}).Error)
	require.NoError(t, dataflowTester.Db.Create(&cqIssueCodeBlockBeforeText{
		DomainEntity: domainlayer.DomainEntity{Id: "existing-domain-block"},
		IssueKey:     "existing-domain-issue",
		Component:    existingComponent,
	}).Error)

	basicRes := implcontext.NewDefaultBasicRes(dataflowTester.Cfg, dataflowTester.Log, dataflowTester.Dal)
	runMigration(t, coremigrations.All(), "change cq_issue_code_blocks.component type to text", basicRes)
	runMigration(t, sonarqubemigrations.All(), "change _tool_sonarqube_issue_code_blocks.component type to text", basicRes)
	assertTextColumnWithoutIndex(t, dataflowTester, "cq_issue_code_blocks")
	assertTextColumnWithoutIndex(t, dataflowTester, "_tool_sonarqube_issue_code_blocks")

	var migratedToolBlock sonarqubeIssueCodeBlockBeforeText
	require.NoError(t, dataflowTester.Db.First(&migratedToolBlock, "id = ?", "existing-tool-block").Error)
	require.Equal(t, existingComponent, migratedToolBlock.Component)
	var migratedDomainBlock cqIssueCodeBlockBeforeText
	require.NoError(t, dataflowTester.Db.First(&migratedDomainBlock, "id = ?", "existing-domain-block").Error)
	require.Equal(t, existingComponent, migratedDomainBlock.Component)
	require.NoError(t, dataflowTester.Db.Delete(&migratedToolBlock).Error)
	require.NoError(t, dataflowTester.Db.Delete(&migratedDomainBlock).Error)

	longComponent256 := "project:" + strings.Repeat("a", 256)
	longComponent500 := "project:" + strings.Repeat("b", 500)
	require.Greater(t, len(longComponent256), 256)
	require.Greater(t, len(longComponent500), 500)

	issueKey := "TEST-LONG-COMPONENT-ISSUE"
	projectKey := "test-long-component-project"
	result := dataflowTester.Db.Create(&models.SonarqubeIssue{
		ConnectionId: 1,
		IssueKey:     issueKey,
		ProjectKey:   projectKey,
		Component:    longComponent500,
		Rule:         "java:S3776",
		Severity:     "CRITICAL",
	})
	require.NoError(t, result.Error)

	codeBlocks := []*models.SonarqubeIssueCodeBlock{
		{
			ConnectionId: 1,
			Id:           "test-long-component-block-256",
			IssueKey:     issueKey,
			Component:    longComponent256,
			Msg:          "component longer than 256 characters",
		},
		{
			ConnectionId: 1,
			Id:           "test-long-component-block-500",
			IssueKey:     issueKey,
			Component:    longComponent500,
			Msg:          "component longer than 500 characters",
		},
	}
	for _, block := range codeBlocks {
		require.NoError(t, dataflowTester.Db.Create(block).Error)
	}

	dataflowTester.Subtask(tasks.ConvertIssueCodeBlocksMeta, &tasks.SonarqubeTaskData{
		Options: &tasks.SonarqubeOptions{
			ConnectionId: 1,
			ProjectKey:   projectKey,
		},
		TaskStartTime: time.Now(),
	})

	var domainBlocks []codequality.CqIssueCodeBlock
	require.NoError(t, dataflowTester.Db.Find(&domainBlocks).Error)
	require.Len(t, domainBlocks, 2)
	require.ElementsMatch(t,
		[]string{longComponent256, longComponent500},
		[]string{domainBlocks[0].Component, domainBlocks[1].Component},
	)

	var toolBlocks []models.SonarqubeIssueCodeBlock
	require.NoError(t, dataflowTester.Db.Where(
		"connection_id = ? AND issue_key = ?", 1, issueKey,
	).Find(&toolBlocks).Error)
	require.Len(t, toolBlocks, 2)
	require.ElementsMatch(t,
		[]string{longComponent256, longComponent500},
		[]string{toolBlocks[0].Component, toolBlocks[1].Component},
	)
}

func runMigration(t *testing.T, scripts []plugin.MigrationScript, name string, basicRes *implcontext.DefaultBasicRes) {
	t.Helper()
	for _, script := range scripts {
		if script.Name() == name {
			require.NoError(t, script.Up(basicRes))
			return
		}
	}
	require.Fail(t, "migration is not registered", name)
}

func assertTextColumnWithoutIndex(t *testing.T, dataflowTester *e2ehelper.DataFlowTester, table string) {
	t.Helper()
	columnTypes, err := dataflowTester.Db.Migrator().ColumnTypes(table)
	require.NoError(t, err)
	for _, columnType := range columnTypes {
		if columnType.Name() == "component" {
			require.Contains(t, strings.ToLower(columnType.DatabaseTypeName()), "text")
			require.False(t, dataflowTester.Db.Migrator().HasIndex(table, "idx_"+table+"_component"))
			return
		}
	}
	require.Fail(t, "component column not found", table)
}
