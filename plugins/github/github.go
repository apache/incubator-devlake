package main // must be main for plugin entry point

import (
	"github.com/merico-dev/lake/logger" // A pseudo type for Plugin Interface implementation
	"github.com/merico-dev/lake/plugins/github/tasks"
)

type Github string

func (plugin Github) Description() string {
	return "To collect and enrich data from GitHub"
}

func (plugin Github) Execute(options map[string]interface{}, progress chan<- float32) {
	logger.Print("start github plugin execution")

	owner, ok := options["owner"]
	if !ok {
		logger.Print("owner is required for GitHub execution")
		return
	}
	ownerString := owner.(string)

	repositoryName, ok := options["repositoryName"]
	if !ok {
		logger.Print("repositoryName is required for GitHub execution")
		return
	}
	repositoryNameString := repositoryName.(string)

	repoId, collectRepoErr := tasks.CollectRepository(ownerString, repositoryNameString)
	if collectRepoErr != nil {
		logger.Error("Could not collect repositories: ", collectRepoErr)
		return
	}

	collectCommitsErr := tasks.CollectCommits(ownerString, repositoryNameString, repoId)
	if collectCommitsErr != nil {
		logger.Error("Could not collect commits: ", collectCommitsErr)
		return
	}
	tasks.CollectChildrenOnCommits(ownerString, repositoryNameString, repoId)

	collectPRsErr := tasks.CollectPullRequests(ownerString, repositoryNameString, repoId)
	if collectPRsErr != nil {
		logger.Error("Could not collect pull requests: ", collectPRsErr)
		return
	}
	tasks.CollectChildrenOnPullRequests(ownerString, repositoryNameString, repoId)

	progress <- 1

	close(progress)

}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Github //nolint
