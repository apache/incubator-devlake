package main // must be main for plugin entry point

import (
	"context"
	"fmt"
	"strings"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger" // A pseudo type for Plugin Interface implementation
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/tasks"
	"github.com/merico-dev/lake/utils"
	"github.com/mitchellh/mapstructure"
)

type GithubOptions struct {
	Tasks []string
}
type Github string

func (plugin Github) Description() string {
	return "To collect and enrich data from GitHub"
}

func (plugin Github) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	configTokensString := config.V.GetString("GITHUB_AUTH")
	tokens := strings.Split(configTokensString, ",")
	githubApiClient := tasks.CreateApiClient(tokens)

	// GitHub API has very low rate limits, so we cycle through multiple tokens to increase rate limits.
	// Then we set the ants max worker per second according to the number of tokens we have to maximize speed.
	tokenCount := len(tokens)
	// We need this rate limit set to 1 by default since there are only 5000 requests per hour allowed for the github api
	maxWorkerPerSecond := 1
	if tokenCount > 0 {
		maxWorkerPerSecond = tokenCount
	}
	scheduler, err := utils.NewWorkerScheduler(50, maxWorkerPerSecond, ctx)
	if err != nil {
		logger.Error("could not create scheduler", false)
	}

	defer scheduler.Release()

	logger.Print("start github plugin execution")

	owner, ok := options["owner"]
	if !ok {
		return fmt.Errorf("owner is required for GitHub execution")
	}
	ownerString := owner.(string)

	repositoryName, ok := options["repositoryName"]
	if !ok {
		return fmt.Errorf("repositoryName is required for GitHub execution")
	}
	repositoryNameString := repositoryName.(string)

	var op GithubOptions
	err = mapstructure.Decode(options, &op)
	if err != nil {
		return err
	}
	tasksToRun := make(map[string]bool, len(op.Tasks))
	for _, task := range op.Tasks {
		tasksToRun[task] = true
	}
	if len(tasksToRun) == 0 {
		tasksToRun = map[string]bool{
			"collectCommits": true,
			"collectIssues":  true,
		}
	}
	repoId, collectRepoErr := tasks.CollectRepository(ownerString, repositoryNameString, githubApiClient)
	if collectRepoErr != nil {
		return fmt.Errorf("Could not collect repositories: %v", collectRepoErr)
	}
	err = tasks.CollectRepositoryIssueLabels(ownerString, repositoryNameString, scheduler, githubApiClient)
	if err != nil {
		return fmt.Errorf("Could not collect repo Issue Labels: %v", err)
	}

	progress <- 0.2

	if tasksToRun["collectCommits"] {
		fmt.Println("INFO >>> starting commits collection")
		collectCommitsErr := tasks.CollectCommits(ownerString, repositoryNameString, repoId, scheduler, githubApiClient)
		if collectCommitsErr != nil {
			return fmt.Errorf("Could not collect commits: %v", collectCommitsErr)
		}
		tasks.CollectChildrenOnCommits(ownerString, repositoryNameString, repoId, scheduler, githubApiClient)
	}

	progress <- 0.4

	if tasksToRun["collectIssues"] {
		fmt.Println("INFO >>> starting issues / PRs collection")
		collectIssuesErr := tasks.CollectIssues(ownerString, repositoryNameString, repoId, scheduler, githubApiClient)
		if collectIssuesErr != nil {
			return fmt.Errorf("Could not collect issues: %v", collectIssuesErr)
		}
		progress <- 0.6

		fmt.Println("INFO >>> starting children on issues collection")
		collectIssueChildrenErr := tasks.CollectChildrenOnIssues(ownerString, repositoryNameString, repoId, scheduler, githubApiClient)
		if collectIssueChildrenErr != nil {
			return fmt.Errorf("Could not collect Issue children: %v", collectIssueChildrenErr)
		}

		progress <- 0.8

		fmt.Println("INFO >>> collecting PR children collection")
		collectPrChildrenErr := tasks.CollectChildrenOnPullRequests(ownerString, repositoryNameString, repoId, scheduler, githubApiClient)
		if collectPrChildrenErr != nil {
			return fmt.Errorf("Could not collect PR children: %v", collectPrChildrenErr)
		}

	}

	if tasksToRun["enrichIssues"] {
		fmt.Println("INFO >>> Enriching Issues")
		enrichmentError := tasks.EnrichIssues(ownerString, repositoryNameString, repoId, scheduler, githubApiClient)
		if enrichmentError != nil {
			return fmt.Errorf("Could not enrich issues: %v", enrichmentError)
		}

	}

	progress <- 1

	close(progress)

	return nil
}

func (plugin Github) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/github"
}

func (plugin Github) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return make(map[string]map[string]core.ApiResourceHandler)
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Github //nolint
