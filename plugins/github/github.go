package main // must be main for plugin entry point

import (
	"context"
	"fmt"

	"github.com/merico-dev/lake/logger" // A pseudo type for Plugin Interface implementation
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/tasks"
	"github.com/merico-dev/lake/utils"
)

type Github string

func (plugin Github) Description() string {
	return "To collect and enrich data from GitHub"
}

func (plugin Github) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	// We need this rate limit set to 1 since there are only 5000 requests per hour allowed for the github api
	scheduler, err := utils.NewWorkerScheduler(50, 1, ctx)
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

	repoId, collectRepoErr := tasks.CollectRepository(ownerString, repositoryNameString)
	if collectRepoErr != nil {
		return fmt.Errorf("Could not collect repositories: %v", collectRepoErr)
	}
	err = tasks.CollectRepositoryIssueLabels(ownerString, repositoryNameString, scheduler)
	if err != nil {
		return fmt.Errorf("Could not collect repo Issue Labels: %v", err)
	}
	progress <- 0.2
	fmt.Println("INFO >>> starting commits collection")
	collectCommitsErr := tasks.CollectCommits(ownerString, repositoryNameString, repoId, scheduler)
	if collectCommitsErr != nil {
		return fmt.Errorf("Could not collect commits: %v", collectCommitsErr)
	}
	tasks.CollectChildrenOnCommits(ownerString, repositoryNameString, repoId, scheduler)
	progress <- 0.4
	fmt.Println("INFO >>> starting issues / PRs collection")
	collectIssuesErr := tasks.CollectIssues(ownerString, repositoryNameString, repoId, scheduler)
	if collectIssuesErr != nil {
		return fmt.Errorf("Could not collect issues: %v", collectIssuesErr)
	}
	progress <- 0.6
	fmt.Println("INFO >>> starting children on issues collection")
	collectIssueChildrenErr := tasks.CollectChildrenOnIssues(ownerString, repositoryNameString, repoId, scheduler)
	if collectIssueChildrenErr != nil {
		return fmt.Errorf("Could not collect Issue children: %v", collectIssueChildrenErr)
	}
	progress <- 0.8
	fmt.Println("INFO >>> collecting PR children collection")
	collectPrChildrenErr := tasks.CollectChildrenOnPullRequests(ownerString, repositoryNameString, repoId, scheduler)
	if collectPrChildrenErr != nil {
		return fmt.Errorf("Could not collect PR children: %v", collectPrChildrenErr)
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
