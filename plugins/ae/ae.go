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
	"github.com/apache/incubator-devlake/plugins/ae/impl"
	"github.com/apache/incubator-devlake/runner"
	"github.com/spf13/cobra"
)

// PluginEntry is a variable exported for Framework to search and load
var PluginEntry impl.AE //nolint

func main() {
	aeCmd := &cobra.Command{Use: "ae"}
	connectionId := aeCmd.Flags().Uint64P("Connection-id", "c", 0, "ae connection id")
	projectId := aeCmd.Flags().IntP("project-id", "p", 0, "ae project id")
	_ = aeCmd.MarkFlagRequired("project-id")
	aeCmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{
			"connectionId": *connectionId,
			"projectId":    *projectId,
		})
	}
	runner.RunCmd(aeCmd)
}
