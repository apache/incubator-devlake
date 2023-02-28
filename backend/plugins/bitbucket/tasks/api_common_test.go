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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"net/http"
	"reflect"
	"testing"
)

func Test_ignoreSomeHTTPStatus(t *testing.T) {
	type args struct {
		statusCodes []int
		res         *http.Response
	}
	tests := []struct {
		name string
		args args
		want errors.Error
	}{
		{
			"400",
			args{
				statusCodes: []int{http.StatusBadRequest, http.StatusNotFound},
				res:         &http.Response{StatusCode: http.StatusBadRequest},
			},
			api.ErrIgnoreAndContinue,
		},
		{
			"404",
			args{
				statusCodes: []int{http.StatusBadRequest, http.StatusNotFound},
				res:         &http.Response{StatusCode: http.StatusNotFound},
			},
			api.ErrIgnoreAndContinue,
		},
		{
			"200",
			args{
				statusCodes: []int{http.StatusBadRequest, http.StatusNotFound},
				res:         &http.Response{StatusCode: http.StatusOK},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ignoreSomeHTTPStatus(tt.args.statusCodes...)(tt.args.res); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ignoreSomeHTTPStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
