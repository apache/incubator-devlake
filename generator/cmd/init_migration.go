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
	"errors"
	"fmt"
	"github.com/apache/incubator-devlake/generator/util"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/stoewer/go-strcase"
	"os"
	"path/filepath"
	"time"
)

func init() {
	rootCmd.AddCommand(initMigrationCmd)
}

var initMigrationCmd = &cobra.Command{
	Use:   "init-migration [plugin_name]",
	Short: "Init migration for plugin",
	Long: `Init migration for plugin
Type in which plugin do you want init migrations in, then generator will create a init migration in plugins/$plugin_name/models/migrationscripts/ for you.`,
	Run: func(cmd *cobra.Command, args []string) {
		var pluginName string

		// try to get plugin name
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
		migrationPath := filepath.Join(`plugins`, pluginName, `models`, `migrationscripts`)
		_, err := os.Stat(migrationPath)
		if !os.IsNotExist(err) {
			cobra.CheckErr(errors.New(`migrationscripts inited or path read file`))
		}

		// read template
		templates := map[string]string{
			`initSchema.go`:     util.ReadTemplate("generator/template/migrationscripts/init_domain_schemas.go-template"),
			`register.go`:       util.ReadTemplate("generator/template/migrationscripts/register.go-template"),
			`archived/.gitkeep`: ``,
		}

		// create vars
		values := map[string]string{}
		util.GenerateAllFormatVar(values, `plugin_name`, pluginName)
		values[`Date`] = time.Now().Format(`20060102`)
		values = util.DetectExistVars(templates, values)
		println(`vars in template:`, fmt.Sprint(values))

		// write template
		util.ReplaceVarInTemplates(templates, values)
		util.WriteTemplates(migrationPath, templates)
		if modifyExistCode {
			println("Last Step: add some code in plugin to implement Migratable like this:\n" +
				"func (plugin " + strcase.UpperCamelCase(pluginName) + ") MigrationScripts() []migration.Script {\n\treturn migrationscripts.All()\n}")
		}
	},
}
