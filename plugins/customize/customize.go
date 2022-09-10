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

package main // must be main for plugin entry point

import (
	"github.com/apache/incubator-devlake/plugins/customize/impl"
	"github.com/apache/incubator-devlake/plugins/customize/tasks"
	"github.com/apache/incubator-devlake/runner"
	"github.com/spf13/cobra"
)

var PluginEntry impl.Customize //nolint

// standalone mode for debugging
func main() {
	cmd := &cobra.Command{Use: "customize"}
	cmd.Run = func(c *cobra.Command, args []string) {
		runner.DirectRun(c, args, PluginEntry, map[string]interface{}{"transformationRules": []tasks.MappingRules{{
			Table:         "issues",
			RawDataTable:  "_raw_jira_api_issues",
			RawDataParams: "{\"ConnectionId\":1,\"BoardId\":8}",
			Mapping:       map[string]string{"x_test": "fields.timespent"},
		}},
		})
	}
	runner.RunCmd(cmd)
}
