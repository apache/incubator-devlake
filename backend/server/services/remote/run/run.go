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
	"encoding/json"
	"fmt"

	"github.com/apache/incubator-devlake/core/runner"
	"github.com/apache/incubator-devlake/server/services/remote"
	"github.com/apache/incubator-devlake/server/services/remote/bridge"
	"github.com/apache/incubator-devlake/server/services/remote/models"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{Use: "run"}

	pluginPath := cmd.Flags().StringP("path", "p", "", "path of the plugin directory")
	connectionId := cmd.Flags().Uint64P("connectionId", "c", 0, "connection id")
	optionsJSON := cmd.Flags().StringP("options", "o", "{}", "plugin options as a JSON object")
	_ = cmd.MarkFlagRequired("path")
	_ = cmd.MarkFlagRequired("connectionId")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		invoker := bridge.NewCmdInvoker(*pluginPath)

		pluginInfo := models.PluginInfo{}
		err := invoker.Call("plugin-info", bridge.DefaultContext).Get(&pluginInfo)

		if err != nil {
			panic(fmt.Sprintf("Cannot get plugin info: %s", err))
		}

		plugin, err := remote.NewRemotePlugin(&pluginInfo)
		if err != nil {
			panic(fmt.Sprintf("Cannot initialize plugin: %s", err))
		}

		var options map[string]interface{}
		jsonErr := json.Unmarshal([]byte(*optionsJSON), &options)
		if jsonErr != nil {
			panic(fmt.Sprintf("Cannot parse options: %s", jsonErr))
		}

		options["connectionId"] = *connectionId
		runner.DirectRun(cmd, args, plugin, options)
	}
	runner.RunCmd(cmd)
}
