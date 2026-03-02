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

var _ plugin.MigrationScript = (*addAccountIdToS3Slice)(nil)

type addAccountIdToS3Slice struct{}

func (*addAccountIdToS3Slice) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()

	err := db.Exec(`
		ALTER TABLE _tool_q_dev_s3_slices
		ADD COLUMN IF NOT EXISTS account_id VARCHAR(255) DEFAULT NULL
	`)
	if err != nil {
		// Try alternative syntax for databases that don't support IF NOT EXISTS
		_ = db.Exec(`ALTER TABLE _tool_q_dev_s3_slices ADD COLUMN account_id VARCHAR(255) DEFAULT NULL`)
	}

	return nil
}

func (*addAccountIdToS3Slice) Version() uint64 {
	return 20260220000001
}

func (*addAccountIdToS3Slice) Name() string {
	return "add account_id column to _tool_q_dev_s3_slices for auto-constructing S3 prefixes"
}
