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
	"github.com/apache/incubator-devlake/plugins/refdiff/impl"
	"github.com/apache/incubator-devlake/runner"
	"github.com/spf13/cobra"
)

var PluginEntry impl.RefDiff

// standalone mode for debugging
func main() {
	refdiffCmd := &cobra.Command{Use: "refdiff"}
	repoId := refdiffCmd.Flags().StringP("repo-id", "r", "", "repo id")
	newRef := refdiffCmd.Flags().StringP("new-ref", "n", "", "new ref")
	oldRef := refdiffCmd.Flags().StringP("old-ref", "o", "", "old ref")

	tagsPattern := refdiffCmd.Flags().StringP("tags-pattern", "p", "", "tags pattern")
	tagsLimit := refdiffCmd.Flags().IntP("tags-limit", "l", 2, "tags limit")
	tagsOrder := refdiffCmd.Flags().StringP("tags-order", "d", "", "tags order")

	_ = refdiffCmd.MarkFlagRequired("repo-id")
	//_ = refdiffCmd.MarkFlagRequired("new-ref")
	//_ = refdiffCmd.MarkFlagRequired("old-ref")

	refdiffCmd.Run = func(cmd *cobra.Command, args []string) {
		pairs := make([]map[string]string, 0, 1)
		if *newRef == "" && *oldRef == "" {
			if *tagsPattern == "" {
				panic("You must set at least one part of '-p' or '-n -o' for tagsPattern or newRef,oldRef")
			}
		} else {
			pairs = append(pairs, map[string]string{
				"NewRef": *newRef,
				"OldRef": *oldRef,
			})
		}

		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{
			"repoId":      repoId,
			"pairs":       pairs,
			"tagsPattern": *tagsPattern,
			"tagsLimit":   *tagsLimit,
			"tagsOrder":   *tagsOrder,
		})
	}
	runner.RunCmd(refdiffCmd)
}
