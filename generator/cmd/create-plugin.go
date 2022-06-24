package cmd

import (
	"errors"
	"fmt"
	"github.com/apache/incubator-devlake/generator/util"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

func init() {
	rootCmd.AddCommand(createPluginCmd)
}

func pluginNameNotExistValidate(input string) error {
	if input == `` {
		return errors.New("plugin name requite")
	}
	snakeNameReg := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)
	if !snakeNameReg.MatchString(input) {
		return errors.New("plugin name invalid (start with a-z and consist with a-z0-9_)")
	}
	_, err := os.Stat(`plugins/` + input)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return errors.New("plugin exists")
}

func pluginNameExistValidate(input string) error {
	if input == `` {
		return errors.New("plugin name requite")
	}
	_, err := os.Stat(`plugins/` + input)
	return err
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
			Label: "with_api_client (is this plugin will request HTTP APIs?)",
			Items: []string{"Yes", "No"},
		}
		_, withApiClient, err := prompt.Run()
		cobra.CheckErr(err)

		values := map[string]string{}
		templates := map[string]string{}
		if withApiClient == `Yes` {
			prompt := promptui.Prompt{
				Label:   "Endpoint (which website/openapi would you want to request)",
				Default: `https://open.example.cn/api/v1`,
			}
			endpoint, err := prompt.Run()
			cobra.CheckErr(err)

			// read template
			templates = map[string]string{
				`plugin_main.go`:      util.ReadTemplate("generator/template/plugin/plugin_main_with_api_client.go-template"),
				`tasks/api_client.go`: util.ReadTemplate("generator/template/plugin/tasks/api_client.go-template"),
				`tasks/task_data.go`:  util.ReadTemplate("generator/template/plugin/tasks/task_data_with_api_client.go-template"),
			}
			util.GenerateAllFormatVar(values, `plugin_name`, pluginName)
			values[`Endpoint`] = endpoint
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
