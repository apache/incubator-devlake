package config

import (
	"fmt"
	"os"

	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

// Lowcase V for private this. You can use it by call GetConfig.
var v *viper.Viper = nil

func GetConfig() *viper.Viper {
	return v
}

// Set default value for no .env or .env not set it
func setDefaultValue() {
	v.SetDefault("DB_URL", "mysql://merico:merico@mysql:3306/lake?charset=utf8mb4&parseTime=True")
	v.SetDefault("PORT", ":8080")
	v.SetDefault("PLUGIN_DIR", "bin/plugins")
	v.SetDefault("TEMPORAL_TASK_QUEUE", "DEVLAKE_TASK_QUEUE")
	v.SetDefault("GITLAB_ENDPOINT", "https://gitlab.com/api/v4/")
	v.SetDefault("GITHUB_ENDPOINT", "https://api.github.com/")
}

// replaceNewEnvItemInOldContent replace old config to new config in env file content
func replaceNewEnvItemInOldContent(v *viper.Viper, envFileContent string) (error, string) {
	// prepare reg exp
	encodeEnvNameReg := regexp.MustCompile(`[^a-zA-Z0-9]`)
	if encodeEnvNameReg == nil {
		return fmt.Errorf("encodeEnvNameReg err"), ``
	}

	for _, key := range v.AllKeys() {
		envName := strings.ToUpper(key)
		val := v.Get(envName)
		encodeEnvName := encodeEnvNameReg.ReplaceAllStringFunc(envName, func(s string) string {
			return fmt.Sprintf(`\%v`, s)
		})
		envItemReg := regexp.MustCompile(fmt.Sprintf(`(?im)^\s*%v\s*\=.*$`, encodeEnvName))
		if envItemReg == nil {
			return fmt.Errorf("regexp err"), ``
		}
		envFileContent = envItemReg.ReplaceAllStringFunc(envFileContent, func(s string) string {
			switch ret := val.(type) {
			case string:
				ret = strings.Replace(ret, `\`, `\\`, -1)
				ret = strings.Replace(ret, `=`, `\=`, -1)
				ret = strings.Replace(ret, `'`, `\'`, -1)
				ret = strings.Replace(ret, `"`, `\"`, -1)
				return fmt.Sprintf(`%v="%v"`, envName, ret)
			default:
				return fmt.Sprintf(`%v="%v"`, envName, ret)
			}
		})
	}
	return nil, envFileContent
}

// WriteConfig save viper to .env file
func WriteConfig(v *viper.Viper) error {
	return WriteConfigAs(v, `.env`)
}

// WriteConfigAs save viper to custom filename
func WriteConfigAs(v *viper.Viper, filename string) error {
	aferoFile := afero.NewOsFs()
	fmt.Println("Attempting to write configuration to .env file.")
	var configType string

	ext := filepath.Ext(filename)
	if ext != "" {
		configType = ext[1:]
	}
	if configType != "env" && configType != "dotenv" {
		return v.WriteConfigAs(filename)
	}

	// FIXME viper just have setter and have no getter so create new configPermissions and file
	flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	configPermissions := os.FileMode(0644)
	file, err := afero.ReadFile(aferoFile, filename)
	if err != nil {
		return err
	}

	envFileContent := string(file)
	f, err := aferoFile.OpenFile(filename, flags, configPermissions)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, key := range v.AllKeys() {
		envName := strings.ToUpper(key)
		if !strings.Contains(envFileContent, envName) {
			envFileContent = fmt.Sprintf("%s\n%s=", envFileContent, envName)
		}
	}
	err, envFileContent = replaceNewEnvItemInOldContent(v, envFileContent)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(envFileContent); err != nil {
		return err
	}
	return f.Sync()
}

func init() {
	// create the object and load the .env file
	v = viper.New()
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = ".env"
	}
	v.SetConfigFile(envPath)
	err := v.ReadInConfig()
	if err != nil {
		logrus.Warn("Failed to read [.env] file:", err)
	}
	v.AutomaticEnv()

	setDefaultValue()
	// This line is essential for reading and writing
	v.WatchConfig()
}
