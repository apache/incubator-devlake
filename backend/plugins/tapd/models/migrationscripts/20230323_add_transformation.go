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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/tapd/models/migrationscripts/archived"
)

type TapdWorkspace20230323 struct {
	TransformationRuleId uint64 `json:"transformationRuleId,omitempty" mapstructure:"transformationRuleId"`
}

func (TapdWorkspace20230323) TableName() string {
	return "_tool_tapd_workspaces"
}

type addTransformation struct{}

func (*addTransformation) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.DropTables(archived.TapdSubWorkspace{})
	if err != nil {
		return err
	}

	// update all _raw_data_params in tapd tables (delete empty company id)
	rawTables := []string{
		// raw data tables
		`_raw_tapd_api_bug_changelogs`,
		`_raw_tapd_api_bug_commits`,
		`_raw_tapd_api_bug_custom_fields`,
		`_raw_tapd_api_bug_status`,
		`_raw_tapd_api_bug_status_last_steps`,
		`_raw_tapd_api_bugs`,
		`_raw_tapd_api_iterations`,
		`_raw_tapd_api_stories`,
		`_raw_tapd_api_story_categories`,
		`_raw_tapd_api_story_changelogs`,
		`_raw_tapd_api_story_commits`,
		`_raw_tapd_api_story_custom_fields`,
		`_raw_tapd_api_story_status`,
		`_raw_tapd_api_story_status_last_steps`,
		`_raw_tapd_api_sub_workspaces`,
		`_raw_tapd_api_task_changelogs`,
		`_raw_tapd_api_task_commits`,
		`_raw_tapd_api_task_custom_fields`,
		`_raw_tapd_api_tasks`,
		`_raw_tapd_api_users`,
		`_raw_tapd_api_workitem_types`,
		`_raw_tapd_api_worklogs`,
	}
	for _, rawTable := range rawTables {
		if db.HasTable(rawTable) {
			err = db.UpdateColumn(
				rawTable,
				`params`,
				dal.DalClause{Expr: "REPLACE(params, '\"CompanyId\":0,', '')"},
				dal.Where(`params LIKE '%\"CompanyId\":0,%'`),
			)
			if err != nil {
				return err
			}
		}
	}

	allTapdTables := []string{
		// tool layer tables
		`_tool_tapd_accounts`,
		`_tool_tapd_bug_changelog_items`,
		`_tool_tapd_bug_changelogs`,
		`_tool_tapd_bug_commits`,
		`_tool_tapd_bug_custom_fields`,
		`_tool_tapd_bug_labels`,
		`_tool_tapd_bug_statuses`,
		`_tool_tapd_bugs`,
		`_tool_tapd_iteration_bugs`,
		`_tool_tapd_iteration_stories`,
		`_tool_tapd_iteration_tasks`,
		`_tool_tapd_iterations`,
		`_tool_tapd_stories`,
		`_tool_tapd_story_bugs`,
		`_tool_tapd_story_categories`,
		`_tool_tapd_story_changelog_items`,
		`_tool_tapd_story_changelogs`,
		`_tool_tapd_story_commits`,
		`_tool_tapd_story_custom_fields`,
		`_tool_tapd_story_labels`,
	}
	for _, allTapdTable := range allTapdTables {
		if db.HasTable(allTapdTable) {
			err = db.UpdateColumn(
				allTapdTable,
				`_raw_data_params`,
				dal.DalClause{Expr: "REPLACE(_raw_data_params, '\"CompanyId\":0,', '')"},
				dal.Where(`_raw_data_params LIKE '%\"CompanyId\":0,%'`),
			)
			if err != nil {
				return err
			}
		}
	}

	return migrationhelper.AutoMigrateTables(
		basicRes,
		&archived.TapdTransformationRule{},
		&TapdWorkspace20230323{},
	)
}

func (*addTransformation) Version() uint64 {
	return 20230323000003
}

func (*addTransformation) Name() string {
	return "Tapd add transformation, update _raw_data_params"
}
