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
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/generator/util"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

func init() {
	rootCmd.AddCommand(createPluginCmd)
}

func pluginNameNotExistValidate(input string) error {
	if input == `` {
		return errors.Default.New("plugin name require")
	}
	snakeNameReg := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)
	if !snakeNameReg.MatchString(input) {
		return errors.Default.New("plugin name invalid (start with a-z and consist with a-z0-9_)")
	}
	if strings.ToLower(input) == `framework` || strings.ToLower(input) == `core` || strings.ToLower(input) == `helper` {
		return errors.Default.New("plugin name cannot be `framework` or `core` or `helper`")
	}
	_, err := os.Stat(`plugins/` + input)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return errors.Default.New("plugin exists")
}

func pluginNameExistValidate(input string) error {
	if input == `` {
		return errors.Default.New("plugin name require")
	}
	_, err := os.Stat(`plugins/` + input)
	return err
}

func pluginNames(withFramework bool) (pluginItems []string, err error) {
	files, err := ioutil.ReadDir(`plugins`)
	if err != nil {
		return nil, err
	}
	if withFramework {
		pluginItems = append(pluginItems, `framework`)
	}
	for _, file := range files {
		if file.IsDir() {
			pluginItems = append(pluginItems, file.Name())
		}
	}
	return pluginItems, nil
}

var createPluginCmd = &cobra.Command{
	Use:   "create-plugin [plugin_name]",
	Short: "Create a new plugin",
	Long: `Create a new plugin
Type in what the name of plugin is, then generator will create a new plugin in plugins/$plugin_name for you`,
	Run: func(cmd *cobra.Command, args []string) {
		var pluginName string

		// try to get plugin name
		if len(args) > 0 {
			pluginName = args[0]
		}
		err := pluginNameNotExistValidate(pluginName)
		if err != nil {
			prompt := promptui.Prompt{
				Label:    "plugin_name",
				Validate: pluginNameNotExistValidate,
				Default:  pluginName,
			}
			pluginName, err = prompt.Run()
			cobra.CheckErr(err)
			pluginName = strings.ToLower(pluginName)
		}

		prompt := promptui.Select{
			Label: "complete_plugin (Will this plugin request HTTP APIs?)",
			Items: []string{"Yes", "No"},
		}
		_, withApiClient, err := prompt.Run()
		cobra.CheckErr(err)

		values := map[string]string{}
		templates := map[string]string{}
		if withApiClient == `Yes` {
			versionTimestamp := time.Now().Format(`20060102`)
			values[`Date`] = versionTimestamp
			// read template
			templates = map[string]string{
				fmt.Sprintf(`%s.go`, pluginName): util.ReadTemplate("generator/template/plugin/plugin_main_complete_plugin.go-template"),
				`impl/impl.go`:                   util.ReadTemplate("generator/template/plugin/impl/impl_complete_plugin.go-template"),
				`tasks/api_client.go`:            util.ReadTemplate("generator/template/plugin/tasks/api_client.go-template"),
				`tasks/task_data.go`:             util.ReadTemplate("generator/template/plugin/tasks/task_data_complete_plugin.go-template"),
				`api/connection.go`:              util.ReadTemplate("generator/template/plugin/api/connection.go-template"),
				`models/connection.go`:           util.ReadTemplate("generator/template/plugin/models/connection.go-template"),
				fmt.Sprintf("models/migrationscripts/%s_add_init_tables.go", versionTimestamp): util.ReadTemplate("generator/template/migrationscripts/add_init_tables.go-template"),
				`models/migrationscripts/register.go`:                                          util.ReadTemplate("generator/template/migrationscripts/register.go-template"),
				`api/init.go`:                                                                  util.ReadTemplate("generator/template/plugin/api/init.go-template"),
				`api/blueprint.go`:                                                             util.ReadTemplate("generator/template/plugin/api/blueprint.go-template"),
			}
			util.GenerateAllFormatVar(values, `plugin_name`, pluginName)
		} else if withApiClient == `No` {
			// read template
			templates = map[string]string{
				`plugin_main.go`:     util.ReadTemplate("generator/template/plugin/plugin_main.go-template"),
				`tasks/task_data.go`: util.ReadTemplate("generator/template/plugin/tasks/task_data.go-template"),
			}
			util.GenerateAllFormatVar(values, `plugin_name`, pluginName)
		}

		values = util.DetectExistVars(templates, values)
		println(`vars in template:`, fmt.Sprint(values))

		// write template
		util.ReplaceVarInTemplates(templates, values)
		util.WriteTemplates(`plugins/`+pluginName, templates)
	},
}
