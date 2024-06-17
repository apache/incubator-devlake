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

package main

import (
	"fmt"

	"github.com/apache/incubator-devlake/core/runner"
	"github.com/apache/incubator-devlake/plugins/issue_trace/impl"
	"github.com/spf13/cobra"
)

var PluginEntry impl.IssueTrace //nolint

func main() {
	cmd := &cobra.Command{Use: "spe"}
	connectionId := cmd.Flags().Uint64P("connection", "s", 0, "jira connection id")
	boardId := cmd.Flags().Uint64P("board", "b", 0, "jira board id")
	plugin := cmd.Flags().StringP("plugin", "p", "", "plugin")
	_ = cmd.MarkFlagRequired("connection")
	_ = cmd.MarkFlagRequired("board")
	_ = cmd.MarkFlagRequired("plugin")

	cmd.Run = func(c *cobra.Command, args []string) {
		fmt.Println("Direct running")
		runner.DirectRun(c, args, PluginEntry, map[string]interface{}{
			"connectionId": *connectionId,
			"boardId":      *boardId,
			"plugin":       plugin,
		}, "")
	}
	runner.RunCmd(cmd)
}
