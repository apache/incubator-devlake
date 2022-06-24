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

package cmd

import (
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(createE2eRawCmd)
}

var createE2eRawCmd = &cobra.Command{
	Use:   "create-e2e-raw [plugin_name] [raw_table_name] [csv_file_name]",
	Short: "Create _raw_table.csv for e2e test",
	Long: `Create _raw_table.csv for e2e test
Type in what the raw_table is, then generator will export and save in plugins/$plugin_name/e2e/_raw_$raw_name.csv for you.`,
	Run: func(cmd *cobra.Command, args []string) {
		var pluginName string
		//var err error

		// try to get plugin name and extractor name
		if len(args) > 0 {
			pluginName = args[0]
		}
		if pluginName == `` {
			pluginItems, err := pluginNames(false)
			cobra.CheckErr(err)
			prompt := promptui.Select{
				Label: "plugin_name",
				Items: pluginItems,
			}
			_, pluginName, err = prompt.Run()
			cobra.CheckErr(err)
		}

		if len(args) > 0 {
			pluginName = args[0]
		}
		if pluginName == `` {
			pluginItems, err := pluginNames(false)
			cobra.CheckErr(err)
			prompt := promptui.Select{
				Label: "plugin_name",
				Items: pluginItems,
			}
			_, pluginName, err = prompt.Run()
			cobra.CheckErr(err)
		}

		rawTablesPath := filepath.Join(`plugins`, pluginName, `e2e`, `raw_tables`)
		dataflowTester := e2ehelper.NewDataFlowTester(nil, "gitlab", nil)
		dataflowTester.ExportRawTable(
			`_raw_github_api_issues`,
			filepath.Join(rawTablesPath, `_raw_github_api_issues.csv`),
		)
		println(rawTablesPath, `_raw_github_api_issues.csv`, ` generated`)
	},
}
