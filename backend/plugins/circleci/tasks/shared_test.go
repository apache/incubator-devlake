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
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/stretchr/testify/assert"
)

func makeResponse(statusCode int) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBufferString("")),
		Request:    &http.Request{},
	}
}

func TestIgnoreDeletedOrBrokenBuilds(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       error
	}{
		{"404 returns ErrIgnoreAndContinue", http.StatusNotFound, api.ErrIgnoreAndContinue},
		{"500 returns ErrIgnoreAndContinue", http.StatusInternalServerError, api.ErrIgnoreAndContinue},
		{"200 returns nil", http.StatusOK, nil},
		{"403 returns nil", http.StatusForbidden, nil},
		{"502 returns nil", http.StatusBadGateway, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := makeResponse(tt.statusCode)
			got := ignoreDeletedOrBrokenBuilds(res)
			if tt.want == nil {
				assert.Nil(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

