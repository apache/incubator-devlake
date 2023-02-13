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
	"testing"

	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/customize/impl"
	"github.com/apache/incubator-devlake/plugins/customize/models"
	"github.com/apache/incubator-devlake/plugins/customize/service"
)

func TestCustomizedFieldDataFlow(t *testing.T) {
	var plugin impl.Customize
	dataflowTester := e2ehelper.NewDataFlowTester(t, "customize", plugin)
	dataflowTester.FlushTabler(&models.CustomizedField{})
	dataflowTester.FlushTabler(&ticket.Issue{})
	svc := service.NewService(dataflowTester.Dal)
	err := svc.CreateField(&models.CustomizedField{
		TbName:      "issues",
		ColumnName:  "x_varchar",
		DisplayName: "test column x_varchar",
		DataType:    "varchar(255)",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = svc.CreateField(&models.CustomizedField{
		TbName:      "issues",
		ColumnName:  "x_text",
		DisplayName: "test column x_text",
		DataType:    "text",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = svc.CreateField(&models.CustomizedField{
		TbName:      "issues",
		ColumnName:  "x_int",
		DisplayName: "test column x_int",
		DataType:    "bigint",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = svc.CreateField(&models.CustomizedField{
		TbName:      "issues",
		ColumnName:  "x_float",
		DisplayName: "test column x_float",
		DataType:    "float",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = svc.CreateField(&models.CustomizedField{
		TbName:      "issues",
		ColumnName:  "x_time",
		DisplayName: "test column x_time",
		DataType:    "timestamp",
	})
	if err != nil {
		t.Fatal(err)
	}
	ff, err := svc.GetFields("issues")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range ff {
		t.Logf("%+v\n", f)
	}
	err = svc.DeleteField("issues", "x_varchar")
	if err != nil {
		t.Fatal(err)
	}
}
