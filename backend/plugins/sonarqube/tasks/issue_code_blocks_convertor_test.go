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

package tasks

import (
	"reflect"
	"testing"

	"github.com/apache/incubator-devlake/core/models/domainlayer/codequality"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
	"github.com/stretchr/testify/require"
)

func TestIssueCodeBlockComponentIsText(t *testing.T) {
	testCases := []struct {
		name  string
		model interface{}
	}{
		{name: "tool layer", model: models.SonarqubeIssueCodeBlock{}},
		{name: "domain layer", model: codequality.CqIssueCodeBlock{}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			field, ok := reflect.TypeOf(testCase.model).FieldByName("Component")
			require.True(t, ok)
			require.Equal(t, "type:text", field.Tag.Get("gorm"))
		})
	}
}
