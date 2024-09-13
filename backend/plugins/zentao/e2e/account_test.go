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
	gocontext "context"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/runner"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/impl"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"github.com/apache/incubator-devlake/plugins/zentao/tasks"
	"testing"
	"time"
)

var basicContext = runner.CreateAppBasicRes()

func getFakeAPIClient() *helper.ApiAsyncClient {
	client, _ := helper.NewApiClient(gocontext.Background(),
		getFakeHomepage()+"/api.php/v1/",
		//"https://zentao.demo.haogs.cn/api.php/v1/",
		nil, time.Second*5, "",
		basicContext,
	)
	return &helper.ApiAsyncClient{
		ApiClient: client,
	}
}

func getFakeHomepage() string {
	return "https://zentaomax.demo.qucheng.cc"
}

func TestZentaoAccountDataFlow(t *testing.T) {

	var zentao impl.Zentao
	dataflowTester := e2ehelper.NewDataFlowTester(t, "zentao", zentao)

	taskData := &tasks.ZentaoTaskData{
		Options: &tasks.ZentaoOptions{
			ConnectionId: 1,
			ProjectId:    3,
		},
		AccountCache: tasks.NewAccountCache(dataflowTester.Dal, 1),
		ApiClient:    getFakeAPIClient(),
		HomePageURL:  getFakeHomepage(),
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_zentao_api_accounts.csv",
		"_raw_zentao_api_accounts")

	// verify extraction
	dataflowTester.FlushTabler(&models.ZentaoAccount{})
	dataflowTester.Subtask(tasks.ExtractAccountMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&models.ZentaoAccount{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/_tool_zentao_accounts.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})

	dataflowTester.FlushTabler(&crossdomain.Account{})
	dataflowTester.Subtask(tasks.ConvertAccountMeta, taskData)
	dataflowTester.VerifyTableWithOptions(&crossdomain.Account{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/users.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}
