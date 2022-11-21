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
	"github.com/apache/incubator-devlake/plugins/github_graphql/impl"
	"github.com/apache/incubator-devlake/runner"
	"github.com/spf13/cobra"
)

// PluginEntry exports a symbol for Framework to load
var PluginEntry impl.GithubGraphql //nolint

// standalone mode for debugging
func main() {
	cmd := &cobra.Command{Use: "githubGraphql"}
	connectionId := cmd.Flags().Uint64P("connectionId", "c", 0, "github connection id")
	owner := cmd.Flags().StringP("owner", "o", "", "github owner")
	repo := cmd.Flags().StringP("repo", "r", "", "github repo")
	_ = cmd.MarkFlagRequired("connectionId")
	_ = cmd.MarkFlagRequired("owner")
	_ = cmd.MarkFlagRequired("repo")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{
			"connectionId": *connectionId,
			"owner":        *owner,
			"repo":         *repo,
		})
	}
	runner.RunCmd(cmd)
}
