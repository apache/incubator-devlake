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
	"github.com/spf13/cobra"

	"github.com/apache/devlake/core/runner"
	"github.com/apache/devlake/plugins/kiro/impl"
)

var PluginEntry impl.Kiro

// standalone mode for debugging
func main() {
	cmd := &cobra.Command{Use: "kiro"}
	connectionId := cmd.Flags().Uint64P("connectionId", "c", 0, "kiro connection id")
	accountId := cmd.Flags().StringP("accountId", "a", "", "AWS account id that Kiro exports for")
	year := cmd.Flags().IntP("year", "y", 0, "year to collect")
	month := cmd.Flags().IntP("month", "m", 0, "month to collect; omit to collect the whole year")

	_ = cmd.MarkFlagRequired("connectionId")
	_ = cmd.MarkFlagRequired("accountId")
	_ = cmd.MarkFlagRequired("year")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		options := map[string]interface{}{
			"connectionId": *connectionId,
			"accountId":    *accountId,
			"year":         *year,
		}
		// A zero month means the whole year, matching the scope model where a
		// nil month widens collection.
		if *month > 0 {
			options["month"] = *month
		}
		runner.DirectRun(cmd, args, PluginEntry, options, "")
	}
	runner.RunCmd(cmd)
}
