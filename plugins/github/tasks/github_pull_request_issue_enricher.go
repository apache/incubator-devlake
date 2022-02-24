package tasks

import (
	"context"
	"github.com/merico-dev/lake/errors"
	"gorm.io/gorm/clause"
	"strings"

	re2 "github.com/dlclark/regexp2"
	"github.com/merico-dev/lake/config"
	lakeModels "github.com/merico-dev/lake/models"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
)

var prBodyCloseRegex *re2.Regexp

func init() {
	var prBodyClose = config.V.GetString("GITHUB_PR_BODY_CLOSE")
	if len(prBodyClose) > 0 {
		prBodyCloseRegex = re2.MustCompile(prBodyClose, 0)
	} else {
		prBodyCloseRegex = re2.MustCompile(`(?mi)(fix|close|resolve|fixes|closes|resolves|fixed|closed|resolved)[\s]*(pingcap\/tidb)?(issue)?( )?((#\d+[ ]?)+)`, 0)

	}
}

func EnrichGithubPullRequestIssue(ctx context.Context, repoId int) (err error) {
	githubPullRequst := &githubModels.GithubPullRequest{}
	cursor, err := lakeModels.Db.Model(&githubPullRequst).
		Where("repo_id = ?", repoId).
		Rows()
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
		err = lakeModels.Db.ScanRows(cursor, githubPullRequst)
		if err != nil {
			return err
		}
		issueNumberListStr := getCloseIssueId(githubPullRequst.Body)

		if issueNumberListStr == "" {
			continue
		}

		issueNumberListStr = strings.TrimPrefix(issueNumberListStr, "#")
		issueNumberList := strings.Split(issueNumberListStr, "#")
		issue := &githubModels.GithubIssue{}

		for _, issueNumber := range issueNumberList {
			issueNumber = strings.TrimSpace(issueNumber)
			err = lakeModels.Db.Where("number = ?", issueNumber).Limit(1).Find(issue).Error
			if err != nil {
				return err
			}
			if issue == nil {
				continue
			}
			githubPullRequstIssue := &githubModels.GithubPullRequestIssue{
				PullRequestId: githubPullRequst.GithubId,
				IssueId:       issue.GithubId,
			}

			err = lakeModels.Db.Clauses(
				clause.OnConflict{UpdateAll: true}).Create(githubPullRequstIssue).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getCloseIssueId(body string) string {
	if prBodyCloseRegex != nil {
		if m, _ := prBodyCloseRegex.FindStringMatch(body); m != nil {
			if len(m.Groups()) > 2 {
				return m.Groups()[len(m.Groups())-2].String()
			}
		}
	}
	return ""
}
