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
	"net/http"
	"testing"

	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/stretchr/testify/assert"
)

func TestIgnoreHTTPStatus404(t *testing.T) {
	cases := []struct {
		name       string
		statusCode int
		wantIgnore bool // expect ErrIgnoreAndContinue (graceful skip, no retry)
		wantErr    bool // expect a real error
	}{
		{"404 no issue tracker -> ignore", http.StatusNotFound, true, false},
		{"410 issue tracker sunset -> ignore", http.StatusGone, true, false},
		{"401 unauthorized -> error", http.StatusUnauthorized, false, true},
		{"200 ok -> continue", http.StatusOK, false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ignoreHTTPStatus404(&http.Response{StatusCode: tc.statusCode})
			switch {
			case tc.wantIgnore:
				assert.Equal(t, api.ErrIgnoreAndContinue, err)
			case tc.wantErr:
				assert.Error(t, err)
			default:
				assert.NoError(t, err)
			}
		})
	}
}
