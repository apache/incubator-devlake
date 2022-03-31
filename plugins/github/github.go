package main // must be main for plugin entry point

import (
	"context"
	"fmt"
	"github.com/merico-dev/lake/errors"
	"os"
	"strings"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger" // A pseudo type for Plugin Interface implementation
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/api"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/github/tasks"
	"github.com/merico-dev/lake/utils"
	"github.com/mitchellh/mapstructure"
)

var _ core.Plugin = (*Github)(nil)

type GithubOptions struct {
	Owner string
	Repo  string
	Tasks []string
}
type Github string

func (plugin Github) Init() {
	logger.Info("INFO >>> init GitHub plugin", true)
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
	)
	if err != nil {
		logger.Error("Error migrating github: ", err)
		panic(err)
	}
}

func (plugin Github) Description() string {
	return "To collect and enrich data from GitHub"
}

func (plugin Github) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	var err error

	// process option from request
	var op GithubOptions
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
	tasksToRun := make(map[string]bool, len(op.Tasks))
	if len(op.Tasks) == 0 {
		tasksToRun = map[string]bool{
			"collectRepo":               true,
			"collectCommits":            true,
			"collectCommitsStat":        true,
			"collectIssues":             true,
			"collectIssueEvents":        true,
			"collectIssueComments":      true,
			"collectPullRequests":       true,
			"collectPullRequestReviews": true,
			"collectPullRequestCommits": true,
			"enrichIssues":              true,
			"enrichPullRequests":        true,
			"enrichComments":            true,
			"enrichPullRequestIssues":   true,
			"convertRepos":              true,
			"convertIssues":             true,
			"convertIssueLabels":        true,
			"convertPullRequests":       true,
			"convertCommits":            true,
			"convertPullRequestCommits": true,
			"convertPullRequestLabels":  true,
			"convertPullRequestIssues":  true,
			"convertNotes":              true,
			"convertUsers":              true,
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
	scheduler, err := utils.NewWorkerScheduler(50, rateLimitPerSecondInt, ctx)
	if err != nil {
		return err
	}
	defer scheduler.Release()
	// TODO: add endpoind, auth validation
	apiClient := tasks.NewGithubApiClient(endpoint, tokens, ctx, scheduler)
	err = apiClient.SetProxy(v.GetString("GITHUB_PROXY"))
	if err != nil {
		return err
	}

	logger.Print("start github plugin execution")

	repoId, err := tasks.CollectRepository(op.Owner, op.Repo, apiClient)
	if err != nil {
		return fmt.Errorf("Could not collect repositories: %v", err)
	}
	if tasksToRun["collectCommits"] {
		progress <- 0.1
		fmt.Println("INFO >>> starting commits collection")
		err = tasks.CollectCommits(op.Owner, op.Repo, repoId, apiClient)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("Could not collect commits: %v", err).Error(),
				SubTaskName: "collectCommits",
			}
		}
	}
	if tasksToRun["collectCommitsStat"] {
		progress <- 0.11
		fmt.Println("INFO >>> starting commits stat collection")
		err = tasks.CollectCommitsStat(op.Owner, op.Repo, repoId, apiClient, rateLimitPerSecondInt, ctx)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("Could not collect commits: %v", err).Error(),
				SubTaskName: "collectCommitsStat",
			}
		}
	}
	if tasksToRun["collectIssues"] {
		progress <- 0.15
		fmt.Println("INFO >>> starting issues collection")
		err = tasks.CollectIssues(op.Owner, op.Repo, repoId, apiClient)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("Could not collect issues: %v", err).Error(),
				SubTaskName: "collectIssues",
			}
		}
	}
	if tasksToRun["collectIssueEvents"] {
		progress <- 0.21
		fmt.Println("INFO >>> starting Issue Events collection")
		err = tasks.CollectIssueEvents(op.Owner, op.Repo, apiClient)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("Could not collect Issue Events: %v", err).Error(),
				SubTaskName: "collectIssueEvents",
			}
		}
	}

	if tasksToRun["collectIssueComments"] {
		progress <- 0.24
		fmt.Println("INFO >>> starting Issue Comments collection")
		err = tasks.CollectIssueComments(op.Owner, op.Repo, apiClient)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("Could not collect Issue Comments: %v", err).Error(),
				SubTaskName: "collectIssueComments",
			}
		}
	}
	if tasksToRun["collectPullRequests"] {
		progress <- 0.28
		fmt.Println("INFO >>> collecting PR collection")
		err = tasks.CollectPullRequests(op.Owner, op.Repo, repoId, apiClient)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("Could not collect PR: %v", err).Error(),
				SubTaskName: "collectPullRequests",
			}
		}
	}

	if tasksToRun["collectPullRequestReviews"] {
		progress <- 0.38
		fmt.Println("INFO >>> collecting PR Reviews collection")
		err = tasks.CollectPullRequestReviews(ctx, op.Owner, op.Repo, repoId, apiClient, rateLimitPerSecondInt)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("Could not collect PR Reviews: %v", err).Error(),
				SubTaskName: "collectPullRequestReviews",
			}
		}
	}

	if tasksToRun["collectPullRequestCommits"] {
		progress <- 0.5
		fmt.Println("INFO >>> starting PR Commits collection")
		err = tasks.CollectPullRequestCommits(ctx, op.Owner, op.Repo, repoId, rateLimitPerSecondInt, apiClient)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("Could not collect PR Commits: %v", err).Error(),
				SubTaskName: "collectPullRequestCommits",
			}
		}
	}

	if tasksToRun["enrichIssues"] {
		progress <- 0.6
		fmt.Println("INFO >>> Enriching Issues")
		err = tasks.EnrichIssues(ctx, repoId)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("could not enrich issues: %v", err).Error(),
				SubTaskName: "enrichIssues",
			}
		}
	}
	if tasksToRun["enrichPullRequests"] {
		progress <- 0.65
		fmt.Println("INFO >>> Enriching PullRequests")
		err = tasks.EnrichPullRequests(ctx, repoId)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("could not enrich PullRequests: %v", err).Error(),
				SubTaskName: "enrichPullRequests",
			}
		}
	}
	if tasksToRun["enrichComments"] {
		progress <- 0.68
		fmt.Println("INFO >>> Enriching comments")
		err = tasks.EnrichComments(ctx, repoId)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("could not enrich PullRequests: %v", err).Error(),
				SubTaskName: "enrichComments",
			}
		}
	}
	if tasksToRun["enrichPullRequestIssues"] {
		progress <- 0.73
		fmt.Println("INFO >>> Enriching PullRequestIssues")
		err = tasks.EnrichPullRequestIssues(ctx, repoId, op.Owner, op.Repo)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("could not enrich PullRequests: %v", err).Error(),
				SubTaskName: "enrichPullRequestIssues",
			}
		}
	}
	if tasksToRun["convertRepos"] {
		progress <- 0.80
		fmt.Println("INFO >>> Converting repos")
		err = tasks.ConvertRepos(ctx)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("could not convert Repos: %v", err).Error(),
				SubTaskName: "convertRepos",
			}
		}
	}
	if tasksToRun["convertIssues"] {
		progress <- 0.85
		fmt.Println("INFO >>> Converting Issues")
		err = tasks.ConvertIssues(ctx, repoId)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("could not convert Issues: %v", err).Error(),
				SubTaskName: "convertIssues",
			}
		}
	}
	if tasksToRun["convertIssueLabels"] {

		progress <- 0.90
		fmt.Println("INFO >>> starting convertIssueLabels")
		err = tasks.ConvertIssueLabels(ctx, repoId)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("could not convert IssueLabels: %v", err).Error(),
				SubTaskName: "convertIssueLabels",
			}
		}
	}
	if tasksToRun["convertPullRequests"] {
		progress <- 0.91
		fmt.Println("INFO >>> starting convertPullRequests")
		err = tasks.ConvertPullRequests(ctx, repoId)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("could not convert PullRequests: %v", err).Error(),
				SubTaskName: "convertPullRequests",
			}
		}
	}
	if tasksToRun["convertPullRequestLabels"] {
		progress <- 0.92
		fmt.Println("INFO >>> starting convertPullRequestLabels")
		err = tasks.ConvertPullRequestLabels(ctx, repoId)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("could not convert PullRequests: %v", err).Error(),
				SubTaskName: "convertPullRequestLabels",
			}
		}
	}
	if tasksToRun["convertCommits"] {
		progress <- 0.93
		fmt.Println("INFO >>> starting convertCommits")
		err = tasks.ConvertCommits(repoId, ctx)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("could not convert Commits: %v", err).Error(),
				SubTaskName: "convertCommits",
			}
		}
	}
	if tasksToRun["convertPullRequestCommits"] {
		progress <- 0.94
		fmt.Println("INFO >>> starting convertPullRequestCommits")
		err = tasks.PrCommitConvertor(ctx, repoId)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("could not convert PullRequestCommits: %v", err).Error(),
				SubTaskName: "convertPullRequestCommits",
			}
		}
	}
	if tasksToRun["convertPullRequestIssues"] {
		progress <- 0.95
		fmt.Println("INFO >>> Converting PullRequestIssues")
		err = tasks.ConvertPullRequestIssues(ctx, repoId)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("could not convert PullRequestCommits: %v", err).Error(),
				SubTaskName: "convertPullRequestIssues",
			}
		}
	}
	if tasksToRun["convertNotes"] {
		progress <- 0.97
		fmt.Println("INFO >>> starting convertNotes")
		err = tasks.ConvertNotes(ctx, repoId)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("could not convert Notes: %v", err).Error(),
				SubTaskName: "convertNotes",
			}
		}
	}
	if tasksToRun["convertUsers"] {
		progress <- 0.98
		fmt.Println("INFO >>> starting convertUsers")
		err = tasks.ConvertUsers(ctx)
		if err != nil {
			return &errors.SubTaskError{
				Message:     fmt.Errorf("could not convert Users: %v", err).Error(),
				SubTaskName: "convertUsers",
			}
		}
	}

	progress <- 1
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
	v := config.GetConfig()
	endpoint := v.GetString("GITHUB_ENDPOINT")
	configTokensString := v.GetString("GITHUB_AUTH")
	tokens := strings.Split(configTokensString, ",")
	githubApiClient := tasks.NewGithubApiClient(endpoint, tokens, nil, nil)
	_ = githubApiClient.SetProxy(v.GetString("GITHUB_PROXY"))
	_, collectRepoErr := tasks.CollectRepository(owner, repo, githubApiClient)
	if collectRepoErr != nil {
		fmt.Println(fmt.Errorf("Could not collect repositories: %v", collectRepoErr))
	}
	go func() {
		err := PluginEntry.Execute(
			map[string]interface{}{
				"owner": owner,
				"repo":  repo,
				"tasks": []string{
					"collectRepo",
					//"collectCommits",
					//"collectCommitsStat",
					//"collectIssues",
					//"collectIssueEvents",
					//"collectIssueComments",
					"collectPullRequests",
					//"collectPullRequestReviews",
					//"collectPullRequestCommits",
					//"enrichIssues",
					//"enrichPullRequests",
					//"enrichComments",
					//"enrichPullRequestIssues",
					//"convertRepos",
					//"convertIssues",
					//"convertIssueLabels",
					"convertPullRequests",
					//"convertCommits",
					//"convertPullRequestCommits",
					//"convertPullRequestLabels",
					//"convertPullRequestIssues",
					//"convertNotes",
					//"convertUsers",
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
