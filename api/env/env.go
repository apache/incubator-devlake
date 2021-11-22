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

func Get(ctx *gin.Context) {
	configJson, err := config.GetConfigJson()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Your json is malformed")
		return
	}
	ctx.JSON(http.StatusOK, configJson)
}

func Set(ctx *gin.Context) {
	var data config.Config
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
	V.Set("JIRA_ISSUE_WORKLOAD_FIELD", data.JIRA_ISSUE_WORKLOAD_FIELD)
	V.Set("JIRA_BOARD_GITLAB_PROJECTS", data.JIRA_BOARD_GITLAB_PROJECTS)
	V.Set("JIRA_ISSUE_BUG_STATUS_MAPPING", data.JIRA_ISSUE_BUG_STATUS_MAPPING)
	V.Set("JIRA_ISSUE_INCIDENT_STATUS_MAPPING", data.JIRA_ISSUE_INCIDENT_STATUS_MAPPING)
	V.Set("JIRA_ISSUE_STORY_STATUS_MAPPING", data.JIRA_ISSUE_STORY_STATUS_MAPPING)
	V.Set("JIRA_ISSUE_TYPE_MAPPING", data.JIRA_ISSUE_TYPE_MAPPING)
	V.Set("GITLAB_ENDPOINT", data.GITLAB_ENDPOINT)
	V.Set("GITLAB_AUTH", data.GITLAB_AUTH)
	V.Set("GITHUB_ENDPOINT", data.GITHUB_ENDPOINT)
	V.Set("GITHUB_AUTH", data.GITHUB_AUTH)
	V.Set("GITHUB_PROXY", data.GITHUB_PROXY)
	V.Set("JENKINS_ENDPOINT", data.JENKINS_ENDPOINT)
	V.Set("JENKINS_USERNAME", data.JENKINS_USERNAME)
	V.Set("JENKINS_PASSWORD", data.JENKINS_PASSWORD)

	err = V.WriteConfig()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Could not write config file")
		return
	}

	ctx.JSON(http.StatusOK, data)
}
