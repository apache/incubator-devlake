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
	"github.com/apache/incubator-devlake/plugins/tapd/impl"
	"github.com/apache/incubator-devlake/runner"
	"github.com/spf13/cobra"
)

var PluginEntry impl.Tapd

// standalone mode for debugging
func main() {
	cmd := &cobra.Command{Use: "tapd"}
	connectionId := cmd.Flags().Uint64P("connection", "c", 0, "tapd connection id")
	workspaceId := cmd.Flags().Uint64P("workspace", "w", 0, "tapd workspace id")
	companyId := cmd.Flags().Uint64P("company", "o", 0, "tapd company id")
	err := cmd.MarkFlagRequired("connection")
	if err != nil {
		panic(err)
	}
	err = cmd.MarkFlagRequired("workspace")
	if err != nil {
		panic(err)
	}

	cmd.Run = func(c *cobra.Command, args []string) {
		runner.DirectRun(c, args, PluginEntry, map[string]interface{}{
			"connectionId": *connectionId,
			"workspaceId":  *workspaceId,
			"companyId":    *companyId,
		})

		// cfg := config.GetConfig()
		// log := logger.Global.Nested(cmd.Use)
		// db, err := runner.NewGormDb(cfg, log)
		// if err != nil {
		// 	panic(err)
		// }
		// wsList := make([]*models.TapdWorkspace, 0)
		// err = db.Find(&wsList, "parent_id = ?", 59169984).Error //nolint TODO: fix the unused err
		// if err != nil {
		// 	panic(err)
		// }
		// for _, v := range wsList {
		// 	*workspaceId = v.ID
		// 	runner.DirectRun(c, args, PluginEntry, map[string]interface{}{
		// 		"connectionId": *connectionId,
		// 		"workspaceId":  *workspaceId,
		// 		"companyId":    *companyId,
		// 	})
		// }
	}
	runner.RunCmd(cmd)
}
