package tasks

import (
	"context"
	"fmt"
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/errors"
	lakeModels "github.com/merico-dev/lake/models"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
	"regexp"
	"strconv"
	"strings"
)

func EnrichPullRequestIssues(ctx context.Context, repoId int, owner string, repo string) (err error) {
	//the pattern before the issue number, sometimes, the issue number is #1098, sometimes it is https://xxx/#1098
	var prBodyCloseRegex *regexp.Regexp
	var prBodyClosePattern string
	var numberPrefix string

	v := config.GetConfig()
	prBodyClosePattern = v.GetString("GITHUB_PR_BODY_CLOSE_PATTERN")
	numberPrefix = v.GetString("GITHUB_PR_BODY_NUMBER_PREFIX")

	numberPattern := fmt.Sprintf(numberPrefix+`\d+[ ]*)+)`, owner, repo)
	if len(prBodyClosePattern) > 0 {
		prPattern := prBodyClosePattern + numberPattern
		prBodyCloseRegex = regexp.MustCompile(prPattern)
	}
	numberPrefixRegex := regexp.MustCompile(numberPrefix)
	charPattern := regexp.MustCompile(`[a-zA-Z\s,]+`)
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
		issueNumberListStr := ""
		//find the matched string in body
		if prBodyCloseRegex != nil {
			issueNumberListStr = prBodyCloseRegex.FindString(githubPullRequst.Body)
		}
		//replace https:// to #, then we can process it later
		if strings.Contains(issueNumberListStr, "https") {
			issueNumberListStr = numberPrefixRegex.ReplaceAllString(issueNumberListStr, "#")
		}
		issueNumberListStr = charPattern.ReplaceAllString(issueNumberListStr, "#")
		//split the string by '#'
		issueNumberList := strings.Split(issueNumberListStr, "#")

		for _, issueNumberStr := range issueNumberList {
			issue := &githubModels.GithubIssue{}
			issueNumberStr = strings.TrimSpace(issueNumberStr)
			//change the issueNumberStr to int, if cannot be changed, just continue
			issueNumber, numFormatErr := strconv.Atoi(issueNumberStr)
			if numFormatErr != nil {
				continue
			}
			err = lakeModels.Db.Where("number = ? and repo_id = ?", issueNumber, repoId).Limit(1).Find(issue).Error
			if err != nil {
				return err
			}
			if issue.Number == 0 {
				continue
			}
			githubPullRequstIssue := &githubModels.GithubPullRequestIssue{
				PullRequestId: githubPullRequst.GithubId,
				IssueId:       issue.GithubId,
				PullNumber:    githubPullRequst.Number,
				IssueNumber:   issue.Number,
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
