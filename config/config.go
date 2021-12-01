package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	PORT                               string `mapstructure:"PORT"`
	DB_URL                             string `mapstructure:"DB_URL"`
	MODE                               string `mapstructure:"MODE"`
	JIRA_ENDPOINT                      string `mapstructure:"JIRA_ENDPOINT"`
	JIRA_BASIC_AUTH_ENCODED            string `mapstructure:"JIRA_BASIC_AUTH_ENCODED"`
	JIRA_ISSUE_EPIC_KEY_FIELD          string `mapstructure:"JIRA_ISSUE_EPIC_KEY_FIELD"`
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
	GITHUB_PROXY                       string `mapstructure:"GITHUB_PROXY"`
	JENKINS_ENDPOINT                   string `mapstructure:"JENKINS_ENDPOINT"`
	JENKINS_USERNAME                   string `mapstructure:"JENKINS_USERNAME"`
	JENKINS_PASSWORD                   string `mapstructure:"JENKINS_PASSWORD"`
	AE_APP_ID                          string `mapstructure:"AE_APP_ID"`
	AE_NONCE_STR                       string `mapstructure:"AE_NONCE_STR"`
	AE_SIGN                            string `mapstructure:"AE_SIGN"`
	AE_ENDPOINT                        string `mapstructure:"AE_ENDPOINT"`
}

var V *viper.Viper

func LoadConfigFile() *viper.Viper {
	V = viper.New()
	V.SetConfigFile(".env")
	_ = V.ReadInConfig()
	V.AutomaticEnv()
	return V
}

func init() {
	V := LoadConfigFile()
	V.SetDefault("PORT", ":8080")
	V.SetDefault("PLUGIN_DIR", "bin/plugins")
	// This line is essential for reading and writing
	V.WatchConfig()
}

func GetConfigJson() (*Config, error) {
	var configJson Config
	err := V.Unmarshal(&configJson)
	if err != nil {
		return nil, err
	}
	return &configJson, nil
}
