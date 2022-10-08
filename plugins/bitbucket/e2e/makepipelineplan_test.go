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

package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/bitbucket/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/impl"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/magiconair/properties/assert"
)

var resp = `
{
  "type": "repository",
  "full_name": "thenicetgp/lake",
  "links": {
    "clone": [
      {
        "name": "https",
        "href": "https://thenicetgp@bitbucket.org/thenicetgp/lake.git"
      },
      {
        "name": "ssh",
        "href": "git@bitbucket.org:thenicetgp/lake.git"
      }
    ]
  }
}
`

func TestMakePipelinePlanCloneURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/repositories/thenicetgp/lake")
		// Send response to be tested
		rw.Write([]byte(resp))
	}))
	// Close the server when test finishes
	defer server.Close()
	var plugin impl.Bitbucket
	dataflowTester := e2ehelper.NewDataFlowTester(t, "bitbucket", plugin)
	connection := &models.BitbucketConnection{
		RestConnection: helper.RestConnection{
			BaseConnection: helper.BaseConnection{
				Model: common.Model{ID: 1},
			},
			Endpoint: server.URL,
		},
		BasicAuth: helper.BasicAuth{
			Username: "thenicetgp",
			Password: "secret",
		},
	}
	api.Init(dataflowTester.Cfg, dataflowTester.Log, dataflowTester.Db)
	err := dataflowTester.Dal.AutoMigrate(connection)
	if err != nil {
		t.Fatal(err)
	}
	defer dataflowTester.FlushTabler(connection)
	err = dataflowTester.Dal.Create(connection)
	if err != nil {
		t.Fatal(err)
	}
	scope := &core.BlueprintScopeV100{
		Entities: []string{core.DOMAIN_TYPE_CODE, core.DOMAIN_TYPE_TICKET, core.DOMAIN_TYPE_CODE_REVIEW, core.DOMAIN_TYPE_CROSS},
		Options: []byte(`{
                            "owner": "thenicetgp",
                            "repo": "lake"
                        }`),
		Transformation: nil,
	}
	plan, err := api.MakePipelinePlan(plugin.SubTaskMetas(), 1, []*core.BlueprintScopeV100{scope})
	if err != nil {
		t.Fatal(err)
	}
	for _, stage := range plan {
		for _, task := range stage {
			if task.Plugin == "gitextractor" {
				assert.Equal(t, task.Options["url"], "https://thenicetgp:secret@bitbucket.org/thenicetgp/lake.git")
			}
		}
	}
}
