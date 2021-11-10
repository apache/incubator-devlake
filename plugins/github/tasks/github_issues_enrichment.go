package tasks

import (
	"fmt"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/utils"
)

func SetPriority(level string) error {
	githubIssue := &models.GithubIssue{}

	// get all high priority
	cursor, err := lakeModels.Db.Model(githubIssue).
		Select("github_issue_labels.*, github_issues.*").
		Joins("join github_issue_label_issues On github_issue_label_issues.issue_id = github_issues.github_id").
		Joins("join github_issue_labels On github_issue_labels.github_id = github_issue_label_issues.issue_label_id").
		Where(fmt.Sprintf("github_issue_labels.name like '%%%v%%'", level)).
		Rows()

	if err != nil {
		return err
	}

	defer cursor.Close()

	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, githubIssue)
		if err != nil {
			return err
		}
		githubIssue.Priority = level
		err = lakeModels.Db.Save(githubIssue).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func EnrichIssues(owner string, repositoryName string, repositoryId int, scheduler *utils.WorkerScheduler, githubApiClient *GithubApiClient) error {
	SetPriority("highest")
	SetPriority("high")
	SetPriority("medium")
	SetPriority("low")
	return nil
}
