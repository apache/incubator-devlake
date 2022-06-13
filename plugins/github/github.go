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
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/runner"
	"github.com/spf13/cobra"
)

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry impl.Github //nolint

// standalone mode for debugging
func main() {
	cmd := &cobra.Command{Use: "github"}
	connectionId := cmd.Flags().Uint64P("connection", "c", 0, "github connection id")
	owner := cmd.Flags().StringP("owner", "o", "", "github owner")
	repo := cmd.Flags().StringP("repo", "r", "", "github repo")
	_ = cmd.MarkFlagRequired("connection")
	_ = cmd.MarkFlagRequired("owner")
	_ = cmd.MarkFlagRequired("repo")

	prType := cmd.Flags().String("prType", "type/(.*)$", "pr type")
	prComponent := cmd.Flags().String("prComponent", "component/(.*)$", "pr component")
	prBodyClosePattern := cmd.Flags().String("prBodyClosePattern", "(?mi)(fix|close|resolve|fixes|closes|resolves|fixed|closed|resolved)[\\s]*.*(((and )?(#|https:\\/\\/github.com\\/%s\\/%s\\/issues\\/)\\d+[ ]*)+)", "pr body close pattern")
	issueSeverity := cmd.Flags().String("issueSeverity", "severity/(.*)$", "issue severity")
	issuePriority := cmd.Flags().String("issuePriority", "^(highest|high|medium|low)$", "issue priority")
	issueComponent := cmd.Flags().String("issueComponent", "component/(.*)$", "issue component")
	issueTypeBug := cmd.Flags().String("issueTypeBug", "^(bug|failure|error)$", "issue type bug")
	issueTypeIncident := cmd.Flags().String("issueTypeIncident", "", "issue type incident")
	issueTypeRequirement := cmd.Flags().String("issueTypeRequirement", "^(feat|feature|proposal|requirement)$", "issue type requirement")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{
			"connectionId":         *connectionId,
			"owner":                *owner,
			"repo":                 *repo,
			"prType":               *prType,
			"prComponent":          *prComponent,
			"prBodyClosePattern":   *prBodyClosePattern,
			"issueSeverity":        *issueSeverity,
			"issuePriority":        *issuePriority,
			"issueComponent":       *issueComponent,
			"issueTypeBug":         *issueTypeBug,
			"issueTypeIncident":    *issueTypeIncident,
			"issueTypeRequirement": *issueTypeRequirement,
		})
	}
	runner.RunCmd(cmd)
}
