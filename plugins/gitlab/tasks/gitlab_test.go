package tasks

import (
	"testing"
)

func TestGitlabPlugin(t *testing.T) {
	projectId := 20103385
	CollectCommits(projectId)
}
