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
	"fmt"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
)

type modifyCustomFieldName struct{}

func (*modifyCustomFieldName) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	var err errors.Error
	issuesNameList := []string{"_tool_tapd_stories", "_tool_tapd_bugs", "_tool_tapd_tasks"}
	for _, issuesName := range issuesNameList {
		switch issuesName {
		case "_tool_tapd_bugs":
			for i := 6; i < 9; i++ {
				err = db.RenameColumn(issuesName, fmt.Sprintf("custom_field%d", i), fmt.Sprintf("custom_field_%d", i))
				if err != nil {
					return err
				}
			}
		case "_tool_tapd_stories":
			fallthrough
		case "_tool_tapd_tasks":
			err = db.RenameColumn(issuesName, `custom_field6`, `custom_field_six`)
			if err != nil {
				return err
			}
			err = db.RenameColumn(issuesName, `custom_field7`, `custom_field_seven`)
			if err != nil {
				return err
			}
			err = db.RenameColumn(issuesName, `custom_field8`, `custom_field_eight`)
			if err != nil {
				return err
			}
		}
		for i := 9; i <= 50; i++ {
			err = db.RenameColumn(issuesName, fmt.Sprintf("custom_field%d", i), fmt.Sprintf("custom_field_%d", i))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (*modifyCustomFieldName) Version() uint64 {
	return 20230411000004
}

func (*modifyCustomFieldName) Name() string {
	return "modify tapd custom field name"
}
