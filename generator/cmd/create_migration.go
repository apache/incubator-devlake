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
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/apache/incubator-devlake/generator/util"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createMigrationCmd)
}

var createMigrationCmd = &cobra.Command{
	Use:   "create-migration [plugin_name/framework]",
	Short: "Create a new migration",
	Long: `Create a new migration
Type in what the purpose of migration is, then generator will create a new migration in plugins/$plugin_name/models/migrationscripts/$date_$purpose.go for you.
If framework passed, generator will create a new migration in models/migrationscripts/$date_$purpose.go`,
	Run: func(cmd *cobra.Command, args []string) {
		var pluginName string
		var purpose string
		var err error

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

		// check if migrationscripts inited
		var migrationPath string
		if pluginName == `framework` {
			migrationPath = filepath.Join(`models`, `migrationscripts`)
		} else {
			migrationPath = filepath.Join(`plugins`, pluginName, `models`, `migrationscripts`)
		}
		_, err = os.Stat(migrationPath)
		if os.IsNotExist(err) {
			cobra.CheckErr(errors.New(`migrationscripts not init. please run init-migration first`))
		}
		cobra.CheckErr(err)

		prompt := promptui.Prompt{
			Label:    "purpose",
			Validate: purposeNotExistValidate,
		}
		purpose, err = prompt.Run()
		cobra.CheckErr(err)

		selector := promptui.Select{
			Label: "with_config (is this migrations will use config?)",
			Items: []string{"No", "Yes"},
		}
		_, withConfig, err := selector.Run()
		cobra.CheckErr(err)

		// create vars
		values := map[string]string{}
		values[`Date`] = time.Now().Format(`20060102`)
		values[`Purpose`] = purpose
		existMigrations, err := ioutil.ReadDir(migrationPath)
		cobra.CheckErr(err)
		values[`Count`] = fmt.Sprintf(`%06d`, len(existMigrations))

		// read template
		templates := map[string]string{}
		if withConfig == `Yes` {
			templates[values[`Date`]+`_`+values[`Purpose`]+`.go`] = util.ReadTemplate("generator/template/migrationscripts/migration_with_config.go-template")
		} else {
			templates[values[`Date`]+`_`+values[`Purpose`]+`.go`] = util.ReadTemplate("generator/template/migrationscripts/migration.go-template")
		}
		values = util.DetectExistVars(templates, values)
		println(`vars in template:`, fmt.Sprint(values))

		// write template
		util.ReplaceVarInTemplates(templates, values)
		util.WriteTemplates(migrationPath, templates)
		if modifyExistCode {
			util.ReplaceVarInFile(
				filepath.Join(migrationPath, `register.go`),
				regexp.MustCompile(`(return +\[]migration\.Script ?\{ ?\n?)((\s*[\w.()]+,\n?)*)(\s*})`),
				fmt.Sprintf("$1$2\t\tnew(%s),\n$4", values[`Purpose`]),
			)
		}
	},
}

func purposeNotExistValidate(input string) error {
	if input == `` {
		return errors.New("purpose require")
	}
	snakeNameReg := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)
	if !snakeNameReg.MatchString(input) {
		return errors.New("purpose invalid (start with a-z and consist with a-z0-9_)")
	}

	return nil
}
