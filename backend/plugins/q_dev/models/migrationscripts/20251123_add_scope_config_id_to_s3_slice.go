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
)

type addScopeConfigIdToS3Slice struct{}

func (*addScopeConfigIdToS3Slice) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()

	// Add scope_config_id column to _tool_q_dev_s3_slices table
	err := db.Exec(`
		ALTER TABLE _tool_q_dev_s3_slices
		ADD COLUMN scope_config_id BIGINT UNSIGNED DEFAULT 0
	`)
	if err != nil {
		return errors.Convert(err)
	}

	return nil
}

func (*addScopeConfigIdToS3Slice) Version() uint64 {
	return 20251123000001
}

func (*addScopeConfigIdToS3Slice) Name() string {
	return "Add scope_config_id column to S3 slice table"
}
