package util

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stoewer/go-strcase"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GenerateAllFormatVar(values map[string]string, baseVarName, baseValue string) {
	values[strcase.LowerCamelCase(baseVarName)] = strcase.LowerCamelCase(baseValue)
	values[strcase.UpperCamelCase(baseVarName)] = strcase.UpperCamelCase(baseValue)
	values[strcase.SnakeCase(baseVarName)] = strcase.SnakeCase(baseValue)
	values[strcase.UpperSnakeCase(baseVarName)] = strcase.UpperSnakeCase(baseValue)
	values[strcase.KebabCase(baseVarName)] = strcase.KebabCase(baseValue)
	values[strcase.UpperKebabCase(baseVarName)] = strcase.UpperKebabCase(baseValue)
}

func ReadTemplate(templateFile string) string {
	f, err := ioutil.ReadFile(templateFile)
	cobra.CheckErr(err)
	return string(f)
}

func WriteTemplates(path string, templates map[string]string) {
	err := os.MkdirAll(path, 0777)
	cobra.CheckErr(err)
	for name, template := range templates {
		err := os.MkdirAll(filepath.Dir(filepath.Join(path, name)), 0777)
		cobra.CheckErr(err)
		err = ioutil.WriteFile(filepath.Join(path, name), []byte(template), 0777)
		cobra.CheckErr(err)
		println(filepath.Join(path, name), ` generated`)
	}
}

func ReplaceVarInFile(filename string, reg *regexp.Regexp, new string) {
	f, err := ioutil.ReadFile(filename)
	cobra.CheckErr(err)
	f = reg.ReplaceAll(f, []byte(new))

	err = ioutil.WriteFile(filename, f, 0777)
	cobra.CheckErr(err)
	println(filename, ` updated`)
}

func DetectExistVars(templates map[string]string, values map[string]string) (newValues map[string]string) {
	newValues = map[string]string{}
	for varName, value := range values {
		for _, template := range templates {
			if strings.Index(template, varName) != -1 {
				newValues[varName] = value
			}
		}
	}
	return newValues
}

func ReplaceVarInTemplates(templates map[string]string, valueMap map[string]string) {
	for i, template := range templates {
		templates[i] = ReplaceVars(template, valueMap)
	}
}

func ReplaceVars(s string, valueMap map[string]string) string {
	for varName, value := range valueMap {
		s = ReplaceVar(s, varName, value)
	}
	return s
}

func ReplaceVar(s, varName, value string) string {
	return strings.ReplaceAll(s, fmt.Sprintf(`{{ .%s }}`, varName), value)
}
