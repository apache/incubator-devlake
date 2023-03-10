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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

func TestMakeTransformationRules(t *testing.T) {
	type args struct {
		rule models.JiraTransformationRule
	}
	tests := []struct {
		name  string
		args  args
		want  *JiraTransformationRule
		want1 errors.Error
	}{
		{"non-null RemotelinkRepoPattern",
			args{rule: models.JiraTransformationRule{
				Name:                       "name",
				EpicKeyField:               "epic",
				StoryPointField:            "story",
				RemotelinkCommitShaPattern: "commit sha pattern",
				RemotelinkRepoPattern:      []byte(`["abc","efg"]`),
				TypeMappings:               []byte(`{"10040":{"standardType":"Incident","statusMappings":null}}`),
			}},
			&JiraTransformationRule{
				Name:                       "name",
				EpicKeyField:               "epic",
				StoryPointField:            "story",
				RemotelinkCommitShaPattern: "commit sha pattern",
				RemotelinkRepoPattern:      []string{"abc", "efg"},
				TypeMappings: map[string]TypeMapping{"10040": {
					StandardType:   "Incident",
					StatusMappings: nil,
				}},
			},
			nil,
		},

		{"null RemotelinkRepoPattern",
			args{rule: models.JiraTransformationRule{
				RemotelinkRepoPattern: nil,
				TypeMappings:          nil,
			}},
			&JiraTransformationRule{
				RemotelinkRepoPattern: nil,
				TypeMappings:          nil,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := MakeTransformationRules(tt.args.rule)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeTransformationRules() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("MakeTransformationRules() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
