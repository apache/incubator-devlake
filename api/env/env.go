package env

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/spf13/viper"
)

var V *viper.Viper

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
	GITHUB_AUTH_TOKENS                 string `mapstructure:"GITHUB_ENDPOINT"`
	GITHUB_ENDPOINT                    string `mapstructure:"GITHUB_ENDPOINT"`
	GITHUB_AUTH                        string `mapstructure:"GITHUB_AUTH"`
	JENKINS_ENDPOINT                   string `mapstructure:"JENKINS_ENDPOINT"`
	JENKINS_USERNAME                   string `mapstructure:"JENKINS_USERNAME"`
	JENKINS_PASSWORD                   string `mapstructure:"JENKINS_PASSWORD"`
}

func Get(ctx *gin.Context) {

	V := config.LoadConfigFile()

	var configJson Config
	V.Unmarshal(&configJson)
	logger.Info("JON >>> config", configJson)
	ctx.JSON(http.StatusOK, configJson)
}

func Set(ctx *gin.Context) {
	var data Config
	err := ctx.MustBindWith(&data, binding.JSON)
	if err != nil {
		logger.Error("", err)
		ctx.JSON(http.StatusBadRequest, "Your json is malformed")
		return
	}

	V := config.LoadConfigFile()

	V.Set("PORT", data.PORT)
	V.Set("DB_URL", data.DB_URL)
	V.Set("MODE", data.MODE)
	V.Set("JIRA_ENDPOINT", data.JIRA_ENDPOINT)
	V.Set("JIRA_BASIC_AUTH_ENCODED", data.JIRA_BASIC_AUTH_ENCODED)
	V.Set("JIRA_ISSUE_EPIC_KEY_FIELD", data.JIRA_ISSUE_EPIC_KEY_FIELD)
	V.Set("JIRA_WORKLOAD_COEFFICIENT", data.JIRA_WORKLOAD_COEFFICIENT)
	V.Set("JIRA_ISSUE_WORKLOAD_FIELD", data.JIRA_ISSUE_WORKLOAD_FIELD)
	V.Set("JIRA_BOARD_GITLAB_PROJECTS", data.JIRA_BOARD_GITLAB_PROJECTS)
	V.Set("JIRA_ISSUE_BUG_STATUS_MAPPING", data.JIRA_ISSUE_BUG_STATUS_MAPPING)
	V.Set("JIRA_ISSUE_INCIDENT_STATUS_MAPPING", data.JIRA_ISSUE_INCIDENT_STATUS_MAPPING)
	V.Set("JIRA_ISSUE_STORY_STATUS_MAPPING", data.JIRA_ISSUE_STORY_STATUS_MAPPING)
	V.Set("JIRA_ISSUE_TYPE_MAPPING", data.JIRA_ISSUE_TYPE_MAPPING)
	V.Set("GITLAB_ENDPOINT", data.GITLAB_ENDPOINT)
	V.Set("GITLAB_AUTH", data.GITLAB_AUTH)
	V.Set("GITHUB_AUTH_TOKENS", data.GITHUB_AUTH_TOKENS)
	V.Set("GITHUB_ENDPOINT", data.GITHUB_ENDPOINT)
	V.Set("GITHUB_AUTH", data.GITHUB_AUTH)
	V.Set("JENKINS_ENDPOINT", data.JENKINS_ENDPOINT)
	V.Set("JENKINS_USERNAME", data.JENKINS_USERNAME)
	V.Set("JENKINS_PASSWORD", data.JENKINS_PASSWORD)

	V.WriteConfig()

	ctx.JSON(http.StatusOK, data)
}
