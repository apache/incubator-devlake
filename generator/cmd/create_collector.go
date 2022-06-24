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
	"regexp"
	"strings"
)

func init() {
	rootCmd.AddCommand(createCollectorCmd)
}

func collectorNameNotExistValidateHoc(pluginName string) func(input string) error {
	collectorNameValidate := func(input string) error {
		if input == `` {
			return errors.New("please input which data would you will collect (snake_format)")
		}
		snakeNameReg := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)
		if !snakeNameReg.MatchString(input) {
			return errors.New("collector name invalid (start with a-z and consist with a-z0-9_)")
		}
		_, err := os.Stat(filepath.Join(`plugins`, pluginName, `tasks`, input+`_collector.go`))
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		return errors.New("collector exists")
	}
	return collectorNameValidate
}

func collectorNameExistValidateHoc(pluginName string) func(input string) error {
	collectorNameValidate := func(input string) error {
		if input == `` {
			return errors.New("please input which data would you will collect (snake_format)")
		}
		_, err := os.Stat(filepath.Join(`plugins`, pluginName, `tasks`, input+`_collector.go`))
		return err
	}
	return collectorNameValidate
}

var createCollectorCmd = &cobra.Command{
	Use:   "create-collector [plugin_name] [collector_name]",
	Short: "Create a new collector",
	Long: `Create a new collector
Type in what the name of collector is, then generator will create a new collector in plugins/$plugin_name/tasks/$collector_name for you`,
	Run: func(cmd *cobra.Command, args []string) {
		var pluginName string
		var collectorName string
		var err error

		// try to get plugin name and collector name
		if len(args) > 0 {
			pluginName = args[0]
		}
		prompt := promptui.Prompt{
			Label:    "plugin_name",
			Validate: pluginNameExistValidate,
			Default:  pluginName,
		}
		pluginName, err = prompt.Run()
		cobra.CheckErr(err)
		pluginName = strings.ToLower(pluginName)

		if len(args) > 1 {
			collectorName = args[1]
		}
		prompt = promptui.Prompt{
			Label:    "collector_data_name",
			Validate: collectorNameNotExistValidateHoc(pluginName),
			Default:  collectorName,
		}
		collectorName, err = prompt.Run()
		cobra.CheckErr(err)
		collectorName = strings.ToLower(collectorName)

		// read template
		templates := map[string]string{
			collectorName + `_collector.go`: util.ReadTemplate("generator/template/plugin/tasks/api_collector.go-template"),
		}

		// create vars
		values := map[string]string{}
		util.GenerateAllFormatVar(values, `plugin_name`, pluginName)
		util.GenerateAllFormatVar(values, `collector_data_name`, collectorName)
		collectorDataNameUpperCamel := strcase.UpperCamelCase(collectorName)
		values = util.DetectExistVars(templates, values)
		println(`vars in template:`, fmt.Sprint(values))

		// write template
		util.ReplaceVarInTemplates(templates, values)
		util.WriteTemplates(filepath.Join(`plugins`, pluginName, `tasks`), templates)
		if modifyExistCode {
			util.ReplaceVarInFile(
				filepath.Join(`plugins`, pluginName, `plugin_main.go`),
				regexp.MustCompile(`(return +\[]core\.SubTaskMeta ?\{ ?\n?)((\s*[\w.]+,\n)*)(\s*})`),
				fmt.Sprintf("$1$2\t\ttasks.Collect%sMeta,\n$4", collectorDataNameUpperCamel),
			)
		}
	},
}
