package config

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

type Config struct {
	PORT                               string `mapstructure:"PORT"`
	DB_URL                             string `mapstructure:"DB_URL"`
	MODE                               string `mapstructure:"MODE"`
	JIRA_ENDPOINT                      string `mapstructure:"JIRA_ENDPOINT"`
	JIRA_BASIC_AUTH_ENCODED            string `mapstructure:"JIRA_BASIC_AUTH_ENCODED"`
	JIRA_ISSUE_EPIC_KEY_FIELD          string `mapstructure:"JIRA_ISSUE_EPIC_KEY_FIELD"`
	JIRA_WORKLOAD_COEFFICIENT          string `mapstructure:"JIRA_WORKLOAD_COEFFICIENT"`
	JIRA_ISSUE_WORKLOAD_FIELD          string `mapstructure:"JIRA_ISSUE_WORKLOAD_FIELD"`
	JIRA_BOARD_GITLAB_PROJECTS         string `mapstructure:"JIRA_BOARD_GITLAB_PROJECTS"`
	JIRA_ISSUE_BUG_STATUS_MAPPING      string `mapstructure:"JIRA_ISSUE_BUG_STATUS_MAPPING"`
	JIRA_ISSUE_INCIDENT_STATUS_MAPPING string `mapstructure:"JIRA_ISSUE_INCIDENT_STATUS_MAPPING"`
	JIRA_ISSUE_STORY_STATUS_MAPPING    string `mapstructure:"JIRA_ISSUE_STORY_STATUS_MAPPING"`
	JIRA_ISSUE_TYPE_MAPPING            string `mapstructure:"JIRA_ISSUE_TYPE_MAPPING"`
	GITLAB_ENDPOINT                    string `mapstructure:"GITLAB_ENDPOINT"`
	GITLAB_AUTH                        string `mapstructure:"GITLAB_AUTH"`
	GITHUB_ENDPOINT                    string `mapstructure:"GITHUB_ENDPOINT"`
	GITHUB_AUTH                        string `mapstructure:"GITHUB_AUTH"`
	JENKINS_ENDPOINT                   string `mapstructure:"JENKINS_ENDPOINT"`
	JENKINS_USERNAME                   string `mapstructure:"JENKINS_USERNAME"`
	JENKINS_PASSWORD                   string `mapstructure:"JENKINS_PASSWORD"`
}

var V *viper.Viper

func LoadConfigFile() *viper.Viper {
	V = viper.New()
	configFile := os.Getenv("ENV_FILE")
	if configFile != "" {
		if !path.IsAbs(configFile) {
			panic("Please set ENV_FILE with absolute path. " +
				"Currently it should only be used for go test to load ENVs.")
		}
		V.SetConfigFile(configFile)
		V.Set("WORKING_DIRECTORY", path.Dir(configFile))
	} else {
		V.SetConfigName(".env")
		V.SetConfigType("env")

		V.AddConfigPath(".")
		// For testing in subdirectories 1 level down (ex. ./config)
		V.AddConfigPath("../")
		V.AddConfigPath("conf")
		V.AddConfigPath("etc")

		execPath, execErr := os.Executable()
		if execErr == nil {
			V.AddConfigPath(path.Dir(execPath))
			V.AddConfigPath(path.Join(path.Dir(execPath), "conf"))
			V.AddConfigPath(path.Join(path.Dir(execPath), "etc"))
		}

		wdPath, _ := os.Getwd()
		V.Set("WORKING_DIRECTORY", wdPath)
	}

	_ = V.ReadInConfig()
	V.AutomaticEnv()
	return V
}

func init() {
	V := LoadConfigFile()
	V.SetDefault("PORT", ":8080")
	// This line is essential for reading and writing
	V.WatchConfig()
}

func GetConfigJson() (*Config, error) {
	V := LoadConfigFile()
	var configJson Config
	err := V.Unmarshal(&configJson)
	if err != nil {
		return nil, err
	}
	return &configJson, nil
}
