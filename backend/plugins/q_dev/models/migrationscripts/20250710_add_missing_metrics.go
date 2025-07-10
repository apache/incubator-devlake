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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*addMissingMetrics)(nil)

type addMissingMetrics struct{}

func (*addMissingMetrics) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()

	// Add all missing metrics columns to _tool_q_dev_user_data table
	// All columns are integer type with default value 0
	// Using snake_case column names to match GORM's default naming convention
	_ = db.Exec(`
        ALTER TABLE _tool_q_dev_user_data 
        ADD COLUMN chat_ai_code_lines INT DEFAULT 0,
        ADD COLUMN chat_messages_interacted INT DEFAULT 0,
        ADD COLUMN chat_messages_sent INT DEFAULT 0,
        ADD COLUMN code_fix_acceptance_event_count INT DEFAULT 0,
        ADD COLUMN code_fix_accepted_lines INT DEFAULT 0,
        ADD COLUMN code_fix_generated_lines INT DEFAULT 0,
        ADD COLUMN code_fix_generation_event_count INT DEFAULT 0,
        ADD COLUMN code_review_failed_event_count INT DEFAULT 0,
        ADD COLUMN dev_acceptance_event_count INT DEFAULT 0,
        ADD COLUMN dev_accepted_lines INT DEFAULT 0,
        ADD COLUMN dev_generated_lines INT DEFAULT 0,
        ADD COLUMN dev_generation_event_count INT DEFAULT 0,
        ADD COLUMN doc_generation_accepted_file_updates INT DEFAULT 0,
        ADD COLUMN doc_generation_accepted_files_creations INT DEFAULT 0,
        ADD COLUMN doc_generation_accepted_line_additions INT DEFAULT 0,
        ADD COLUMN doc_generation_accepted_line_updates INT DEFAULT 0,
        ADD COLUMN doc_generation_event_count INT DEFAULT 0,
        ADD COLUMN doc_generation_rejected_file_creations INT DEFAULT 0,
        ADD COLUMN doc_generation_rejected_file_updates INT DEFAULT 0,
        ADD COLUMN doc_generation_rejected_line_additions INT DEFAULT 0,
        ADD COLUMN doc_generation_rejected_line_updates INT DEFAULT 0,
        ADD COLUMN test_generation_accepted_lines INT DEFAULT 0,
        ADD COLUMN test_generation_accepted_tests INT DEFAULT 0,
        ADD COLUMN test_generation_event_count INT DEFAULT 0,
        ADD COLUMN test_generation_generated_lines INT DEFAULT 0,
        ADD COLUMN test_generation_generated_tests INT DEFAULT 0,
        ADD COLUMN transformation_event_count INT DEFAULT 0,
        ADD COLUMN transformation_lines_generated INT DEFAULT 0,
        ADD COLUMN transformation_lines_ingested INT DEFAULT 0
    `)

	return nil
}

func (*addMissingMetrics) Version() uint64 {
	return 20250710000001
}

func (*addMissingMetrics) Name() string {
	return "add missing metrics columns to QDevUserData table"
}
