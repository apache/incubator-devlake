package tasks

import (
	"database/sql"
	"fmt"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/github/models"
)

func buildLabelQuery(matches []string) (*sql.Rows, error) {
	var where string

	for index, s := range matches {
		where += s
		if index != (len(matches) - 1) {
			where += "|"
		}
	}

	cursor, err := lakeModels.Db.Model(&models.GithubIssue{}).
		Select("github_issue_labels.*, github_issues.*").
		Joins("join github_issue_label_issues On github_issue_label_issues.issue_id = github_issues.github_id").
		Joins("join github_issue_labels On github_issue_labels.github_id = github_issue_label_issues.issue_label_id").
		Where(fmt.Sprintf("github_issue_labels.name rlike '.*(%v).*'", where)).
		Rows()

	if err != nil {
		return nil, err
	}

	return cursor, nil
}

func SetPriority(matches []string, value string) error {
	githubIssue := &models.GithubIssue{}

	cursor, _ := buildLabelQuery(matches)
	defer cursor.Close()

	for cursor.Next() {
		err := lakeModels.Db.ScanRows(cursor, githubIssue)
		if err != nil {
			return err
		}
		githubIssue.Priority = value
		err = lakeModels.Db.Save(githubIssue).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func SetType(matches []string, value string) error {
	githubIssue := &models.GithubIssue{}

	cursor, _ := buildLabelQuery(matches)
	defer cursor.Close()

	for cursor.Next() {
		err := lakeModels.Db.ScanRows(cursor, githubIssue)
		if err != nil {
			return err
		}
		githubIssue.Type = value
		err = lakeModels.Db.Save(githubIssue).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func EnrichIssues() error {
	err := SetPriority([]string{"highest"}, "Highest")
	if err != nil {
		return err
	}
	err = SetPriority([]string{"high"}, "High")
	if err != nil {
		return err
	}
	err = SetPriority([]string{"medium"}, "Medium")
	if err != nil {
		return err
	}
	err = SetPriority([]string{"low"}, "Low")
	if err != nil {
		return err
	}

	err = SetType([]string{"bug"}, "Bug")
	if err != nil {
		return err
	}
	err = SetType([]string{"feat", "feature", "proposal"}, "Requirement")
	if err != nil {
		return err
	}
	err = SetType([]string{"doc"}, "Documentation")
	if err != nil {
		return err
	}
	return nil
}
