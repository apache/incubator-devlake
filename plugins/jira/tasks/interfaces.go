package tasks

import (
	"context"
	"github.com/merico-dev/lake/plugins/jira/models"
	"time"
)

type Collector interface {
	CollectBoard(jiraApiClient *JiraApiClient, source *models.JiraSource, boardId uint64) error
	CollectChangelogs(
		jiraApiClient *JiraApiClient,
		source *models.JiraSource,
		boardId uint64,
		rateLimitPerSecondInt int,
		ctx context.Context,
	) error
	CollectIssues(
		jiraApiClient *JiraApiClient,
		source *models.JiraSource,
		boardId uint64,
		since time.Time,
		rateLimitPerSecondInt int,
		ctx context.Context,
	) error
	CollectProjects(jiraApiClient *JiraApiClient, sourceId uint64) error
	CollectRemoteLinks(
		jiraApiClient *JiraApiClient,
		source *models.JiraSource,
		boardId uint64,
		rateLimitPerSecondInt int,
		ctx context.Context,
	) error
	CollectSprint(jiraApiClient *JiraApiClient, source *models.JiraSource, boardId uint64) error
	CollectUsers(jiraApiClient *JiraApiClient, sourceId uint64) error
}
