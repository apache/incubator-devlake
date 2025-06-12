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

type modifyFileMetaTable struct{}

func (*modifyFileMetaTable) Name() string {
	return "Modify QDevS3FileMeta table to allow NULL processed_time"
}

func (*modifyFileMetaTable) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()

	// 修改 processed_time 列允许为 NULL
	sql := "ALTER TABLE _tool_q_dev_s3_file_meta MODIFY processed_time DATETIME NULL"
	err := db.Exec(sql)
	if err != nil {
		return errors.Default.Wrap(err, "failed to modify processed_time column")
	}

	return nil
}

func (*modifyFileMetaTable) Version() uint64 {
	return 20250320
}
