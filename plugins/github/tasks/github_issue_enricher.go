package tasks

import (
	"context"
	"github.com/merico-dev/lake/errors"
	"regexp"

	"github.com/merico-dev/lake/config"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
)

func EnrichIssues(ctx context.Context, repoId int) (err error) {
	var issueSeverityRegex *regexp.Regexp
	var issueComponentRegex *regexp.Regexp
	var issuePriorityRegex *regexp.Regexp
	var issueTypeBugRegex *regexp.Regexp
	var issueTypeRequirementRegex *regexp.Regexp
	var issueTypeIncidentRegex *regexp.Regexp

	v := config.GetConfig()
	var issueSeverity = v.GetString("GITHUB_ISSUE_SEVERITY")
	var issueComponent = v.GetString("GITHUB_ISSUE_COMPONENT")
	var issuePriority = v.GetString("GITHUB_ISSUE_PRIORITY")
	var issueTypeBug = v.GetString("GITHUB_ISSUE_TYPE_BUG")
	var issueTypeRequirement = v.GetString("GITHUB_ISSUE_TYPE_REQUIREMENT")
	var issueTypeIncident = v.GetString("GITHUB_ISSUE_TYPE_INCIDENT")
	if len(issueSeverity) > 0 {
		issueSeverityRegex = regexp.MustCompile(issueSeverity)
	}
	if len(issueComponent) > 0 {
		issueComponentRegex = regexp.MustCompile(issueComponent)
	}
	if len(issuePriority) > 0 {
		issuePriorityRegex = regexp.MustCompile(issuePriority)
	}
	if len(issueTypeBug) > 0 {
		issueTypeBugRegex = regexp.MustCompile(issueTypeBug)
	}
	if len(issueTypeRequirement) > 0 {
		issueTypeRequirementRegex = regexp.MustCompile(issueTypeRequirement)
	}
	if len(issueTypeIncident) > 0 {
		issueTypeIncidentRegex = regexp.MustCompile(issueTypeIncident)
	}

	githubIssue := &githubModels.GithubIssue{}
	cursor, err := lakeModels.Db.Model(&githubIssue).Where("repo_id = ?", repoId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.TaskCanceled
		default:
		}
		err = lakeModels.Db.ScanRows(cursor, githubIssue)
		if err != nil {
			return err
		}
		githubIssue.Severity = ""
		githubIssue.Component = ""
		githubIssue.Priority = ""
		githubIssue.Type = ""

		var issueLabels []string

		err = lakeModels.Db.Table("github_issue_labels").
			Where("issue_id = ?", githubIssue.GithubId).
			Pluck("`label_name`", &issueLabels).Error
		if err != nil {
			return err
		}

		// enrich issues by issue_labels
		for _, issueLabel := range issueLabels {
			if issueSeverityRegex != nil {
				groups := issueSeverityRegex.FindStringSubmatch(issueLabel)
				if len(groups) > 0 {
					githubIssue.Severity = groups[1]
					continue
				}
			}

			if issueComponentRegex != nil {
				groups := issueComponentRegex.FindStringSubmatch(issueLabel)
				if len(groups) > 0 {
					githubIssue.Component = groups[1]
					continue
				}
			}

			if issuePriorityRegex != nil {
				groups := issuePriorityRegex.FindStringSubmatch(issueLabel)
				if len(groups) > 0 {
					githubIssue.Priority = groups[1]
					continue
				}
			}

			if issueTypeBugRegex != nil {
				if ok := issueTypeBugRegex.MatchString(issueLabel); ok {
					githubIssue.Type = ticket.BUG
					continue
				}
			}

			if issueTypeRequirementRegex != nil {
				if ok := issueTypeRequirementRegex.MatchString(issueLabel); ok {
					githubIssue.Type = ticket.REQUIREMENT
					continue
				}
			}

			if issueTypeIncidentRegex != nil {
				if ok := issueTypeIncidentRegex.MatchString(issueLabel); ok {
					githubIssue.Type = ticket.INCIDENT
					continue
				}
			}
		}

		err = lakeModels.Db.Save(githubIssue).Error
		if err != nil {
			return err
		}
	}
	return nil
}

