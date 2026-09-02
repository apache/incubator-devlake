/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License. You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/devlake/plugins/kiro/models"
)

func TestNewConnectionReportSatisfiesConfigUIContract(t *testing.T) {
	connection := &models.KiroConnection{
		KiroConn: models.KiroConn{
			Bucket:          "report-bucket",
			PromptLogBucket: "log-bucket",
		},
	}

	report := newConnectionReport(connection)
	assert.True(t, report.Success)
	assert.Equal(t, "success", report.Message)
	assert.Equal(t, "report-bucket", report.ReportBucket)
	assert.Equal(t, "log-bucket", report.PromptLogBucket)
	assert.NotNil(t, report.Accounts)
	assert.NotNil(t, report.Streams)

	body, err := json.Marshal(report)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"success": true,
		"message": "success",
		"reportBucket": "report-bucket",
		"promptLogBucket": "log-bucket",
		"accounts": [],
		"streams": []
	}`, string(body))
}
