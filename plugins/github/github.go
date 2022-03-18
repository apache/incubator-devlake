package main // must be main for plugin entry point

import (
	"context"
	"fmt"
	lakeModels "github.com/merico-dev/lake/models"
	"os"
	"strings"
	"time"

	"github.com/merico-dev/lake/errors"
	"github.com/merico-dev/lake/plugins/helper"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/api"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/github/tasks"
	"github.com/merico-dev/lake/utils"
	"github.com/mitchellh/mapstructure"
)

var _ core.Plugin = (*Github)(nil)

type Github string

func (plugin Github) Init() {
	err := lakeModels.Db.AutoMigrate(
		&models.GithubRepo{},
		&models.GithubCommit{},
		&models.GithubRepoCommit{},
		&models.GithubPullRequest{},
		&models.GithubReviewer{},
		&models.GithubPullRequestComment{},
		&models.GithubPullRequestCommit{},
		&models.GithubPullRequestLabel{},
		&models.GithubIssue{},
		&models.GithubIssueComment{},
		&models.GithubIssueEvent{},
		&models.GithubIssueLabel{},
		&models.GithubUser{},
		&models.GithubPullRequestIssue{},
		&models.GithubCommitStat{},
	)
	if err != nil {
		panic(err)
	}
}

func (plugin Github) Description() string {
	return "To collect and enrich data from GitHub"
}

func (plugin Github) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	var err error

	// process option from request
	var op tasks.GithubOptions
	err = mapstructure.Decode(options, &op)
	if err != nil {
		return err
	}
	if op.Owner == "" {
		return fmt.Errorf("owner is required for GitHub execution")
	}
	if op.Repo == "" {
		return fmt.Errorf("repo is required for GitHub execution")
	}
	var since *time.Time
	if op.Since != "" {
		*since, err = time.Parse("2006-01-02T15:04:05Z", op.Since)
		if err != nil {
			return fmt.Errorf("invalid value for `since`: %w", err)
		}
	}
	logger := helper.NewDefaultTaskLogger(nil, "github")

	tasksToRun := make(map[string]bool, len(op.Tasks))
	if len(op.Tasks) == 0 {
		tasksToRun = map[string]bool{
			"collectApiRepositories": false,
			"extractApiRepositories": false,
			//"collectApiCommits":
			//"extractApiCommits":
			//"collectApiCommitStats":
			//"extractApiCommitStats":        false,
			"collectIssues":                true,
			"collectIssueEvents":           true,
			"collectIssueComments":         true,
			"collectApiIssues":             false,
			"extractApiIssues":             false,
			"collectApiPullRequests":       false,
			"extractApiPullRequests":       false,
			"collectApiComments":           false,
			"extractApiComments":           false,
			"collectApiEvents":             false,
			"extractApiEvents":             false,
			"collectApiPullRequestCommits": false,
			"extractApiPullRequestCommits": false,
			"collectPullRequests":          true,
			"collectPullRequestReviews":    true,
			"collectPullRequestCommits":    true,
			"enrichIssues":                 true,
			"enrichPullRequests":           true,
			"enrichComments":               true,
			"enrichPullRequestIssues":      true,
			"convertRepos":                 true,
			"convertIssues":                true,
			"convertBoardIssues":           false,
			"convertIssueLabels":           true,
			"convertPullRequests":          true,
			"convertCommits":               true,
			"convertPullRequestCommits":    true,
			"convertPullRequestLabels":     true,
			"convertPullRequestIssues":     true,
			"convertNotes":                 true,
			"convertUsers":                 true,
		}
	} else {
		for _, task := range op.Tasks {
			tasksToRun[task] = true
		}
	}

	v := config.GetConfig()
	// process configuration
	endpoint := v.GetString("GITHUB_ENDPOINT")
	tokens := strings.Split(v.GetString("GITHUB_AUTH"), ",")
	// setup rate limit
	tokenCount := len(tokens)
	if tokenCount == 0 {
		return fmt.Errorf("owner is required for GitHub execution")
	}
	rateLimitPerSecondInt, err := core.GetRateLimitPerSecond(options, tokenCount)
	if err != nil {
		return err
	}
	scheduler, err := utils.NewWorkerScheduler(20, rateLimitPerSecondInt, ctx)
	if err != nil {
		return err
	}
	defer scheduler.Release()
	// TODO: add endpoind, auth validation
	apiClient := tasks.NewGithubApiClient(endpoint, tokens, v.GetString("GITHUB_PROXY"), ctx, scheduler, logger)

	taskData := &tasks.GithubTaskData{
		Options:   &op,
		ApiClient: &apiClient.ApiClient,
		Since:     since,
	}

	taskCtx := helper.NewDefaultTaskContext("github", ctx, logger, taskData, tasksToRun)
	repo := models.GithubRepo{}
	err = taskCtx.GetDb().Model(&models.GithubRepo{}).
		Where("owner_login = ? and name = ?", taskData.Options.Owner, taskData.Options.Repo).Limit(1).Find(&repo).Error
	if err != nil {
		return err
	}

	taskData.Repo = &repo

	newTasks := []struct {
		name       string
		entryPoint core.SubTaskEntryPoint
	}{
		//{name: "collectApiRepositories", entryPoint: tasks.CollectApiRepositories},
		//{name: "extractApiRepositories", entryPoint: tasks.ExtractApiRepositories},
		//{name: "collectApiIssues", entryPoint: tasks.CollectApiIssues},
		//{name: "extractApiIssues", entryPoint: tasks.ExtractApiIssues},
		//{name: "collectApiPullRequests", entryPoint: tasks.CollectApiPullRequests},
		//{name: "extractApiPullRequests", entryPoint: tasks.ExtractApiPullRequests},
		//{name: "collectApiComments", entryPoint: tasks.CollectApiComments},
		//{name: "extractApiComments", entryPoint: tasks.ExtractApiComments},
		//{name: "collectApiEvents", entryPoint: tasks.CollectApiEvents},
		//{name: "extractApiEvents", entryPoint: tasks.ExtractApiEvents},
		//{name: "collectApiPullRequestCommits", entryPoint: tasks.CollectApiPullRequestCommits},
		//{name: "extractApiPullRequestCommits", entryPoint: tasks.ExtractApiPullRequestCommits},
		//{name: "collectApiPullRequestReviews", entryPoint: tasks.CollectApiPullRequestReviews},
		//{name: "extractApiPullRequestReviews", entryPoint: tasks.ExtractApiPullRequestReviews},
		//{name: "collectApiCommits", entryPoint: tasks.CollectApiCommits},
		//{name: "extractApiCommits", entryPoint: tasks.ExtractApiCommits},
		//{name: "collectApiCommitStats", entryPoint: tasks.CollectApiCommitStats},
		//{name: "extractApiCommitStats", entryPoint: tasks.ExtractApiCommitStats},
		{name: "enrichPullRequestIssues", entryPoint: tasks.EnrichPullRequestIssues},
		//{name: "convertRepos", entryPoint: tasks.ConvertRepos},
		//{name: "convertIssues", entryPoint: tasks.ConvertIssues},
		//{name: "convertCommits", entryPoint: tasks.ConvertCommits},
		//{name: "convertIssueLabels", entryPoint: tasks.ConvertIssueLabels},
		//{name: "convertPullRequestCommits", entryPoint: tasks.ConvertPullRequestCommits},
		//{name: "convertPullRequests", entryPoint: tasks.ConvertPullRequests},
		//{name: "convertPullRequestLabels", entryPoint: tasks.ConvertPullRequestLabels},
		{name: "convertPullRequestIssues", entryPoint: tasks.ConvertPullRequestIssues},
		//{name: "convertNotes", entryPoint: tasks.ConvertNotes},
		//{name: "convertUsers", entryPoint: tasks.ConvertUsers},
	}
	for _, t := range newTasks {
		c, err := taskCtx.SubTaskContext(t.name)
		if err != nil {
			return err
		}
		if c != nil {
			err = t.entryPoint(c)
			if err != nil {
				return &errors.SubTaskError{
					SubTaskName: t.name,
					Message:     err.Error(),
				}
			}
		}
	}
	return nil
}

func (plugin Github) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/github"
}

func (plugin Github) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"sources": {
			"GET":  api.ListSources,
			"POST": api.PutSource,
		},
		"sources/:sourceId": {
			"GET": api.GetSource,
			"PUT": api.PutSource,
		},
	}
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Github //nolint

// standalone mode for debugging
func main() {
	args := os.Args[1:]
	owner := "merico-dev"
	repo := "lake"
	if len(args) > 0 {
		owner = args[0]
	}
	if len(args) > 1 {
		repo = args[1]
	}

	err := core.RegisterPlugin("github", PluginEntry)
	if err != nil {
		panic(err)
	}
	PluginEntry.Init()
	progress := make(chan float32)

	go func() {
		err := PluginEntry.Execute(
			map[string]interface{}{
				"owner": owner,
				"repo":  repo,
				"tasks": []string{
					//"collectApiRepositories",
					//"extractApiRepositories",
					//"collectApiCommits",
					//"extractApiCommits",
					//"collectApiCommitStats",
					//"extractApiCommitStats",
					//
					//"collectApiIssues",
					//"extractApiIssues",
					//"collectApiComments",
					//"extractApiComments",
					//"collectApiEvents",
					//"extractApiEvents",
					//"collectApiPullRequests",
					//"extractApiPullRequests",
					//"collectApiPullRequestCommits",
					//"extractApiPullRequestCommits",
					//"collectApiPullRequestReviews",
					//"extractApiPullRequestReviews",
					"enrichPullRequestIssues",
					"convertRepos",
					"convertIssues",
					"convertIssueLabels",
					"convertPullRequests",
					"convertCommits",
					"convertPullRequestCommits",
					"convertPullRequestLabels",
					"convertPullRequestIssues",
					"convertNotes",
					"convertUsers",
				},
			},
			progress,
			context.Background(),
		)
		if err != nil {
			panic(err)
		}
		close(progress)
	}()
	for p := range progress {
		fmt.Println(p)
	}
}
