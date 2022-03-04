package tasks

import (
	"context"
	"github.com/merico-dev/lake/plugins/jira/models"
	"time"
)

type Cloud struct {
}

func NewCloud() *Cloud {
	return &Cloud{}
}

func (c *Cloud) CollectBoard(jiraApiClient *JiraApiClient, source *models.JiraSource, boardId uint64) error {
	return CollectBoard(jiraApiClient, source, boardId)
}

func (c *Cloud) CollectIssues(jiraApiClient *JiraApiClient, source *models.JiraSource, boardId uint64, since time.Time, rateLimitPerSecondInt int, ctx context.Context) error {
	return CollectIssues(jiraApiClient, source, boardId, since, rateLimitPerSecondInt, ctx)
}

func (c *Cloud) CollectProjects(jiraApiClient *JiraApiClient, sourceId uint64) error {
	return CollectProjects(jiraApiClient, sourceId)
}

func (c *Cloud) CollectRemoteLinks(jiraApiClient *JiraApiClient, source *models.JiraSource, boardId uint64, rateLimitPerSecondInt int, ctx context.Context) error {
	return CollectRemoteLinks(jiraApiClient, source, boardId, rateLimitPerSecondInt, ctx, collectRemotelinksByIssueId)
}

func (c *Cloud) CollectChangelogs(
	jiraApiClient *JiraApiClient,
	source *models.JiraSource,
	boardId uint64,
	rateLimitPerSecondInt int,
	ctx context.Context,
) error {
	return CollectChangelogs(jiraApiClient, source, boardId, rateLimitPerSecondInt, ctx, collectChangelogsByIssueId)
}

func (c *Cloud) CollectSprint(jiraApiClient *JiraApiClient, source *models.JiraSource, boardId uint64) error {
	return CollectSprint(jiraApiClient, source, boardId)
}

func (c *Cloud) CollectUsers(jiraApiClient *JiraApiClient, sourceId uint64) error {
	return CollectUsers(jiraApiClient, sourceId)
}
