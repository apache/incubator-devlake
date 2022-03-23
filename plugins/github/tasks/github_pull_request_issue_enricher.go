package tasks

import (
	"context"
	"fmt"
	"github.com/merico-dev/lake/config"
	lakeModels "github.com/merico-dev/lake/models"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
	"regexp"
	"strconv"
	"strings"
)

var prBodyCloseRegex *regexp.Regexp
var prBodyClosePattern string
var numberPrefix string

func init() {
	prBodyClosePattern = config.GetConfig().GetString("GITHUB_PR_BODY_CLOSE_PATTERN")
	numberPrefix = config.GetConfig().GetString("GITHUB_PR_BODY_NUMBER_PREFIX")
}

func EnrichPullRequestIssues(ctx context.Context, repoId int, owner string, repo string) (err error) {
	//the pattern before the issue number, sometimes, the issue number is #1098, sometimes it is https://xxx/#1098
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
	resList := make([]string, 0)

	defer cursor.Close()
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		err = lakeModels.Db.ScanRows(cursor, githubPullRequst)
		if err != nil {
			return err
		}

		//find the matched string in body
		issueNumberListStr := getCloseIssueId(githubPullRequst.Body)

		if issueNumberListStr == "" {
			continue
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
			if issue == nil {
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
	for _, v := range resList {
		fmt.Println(v)
	}

	return nil
}

func getCloseIssueId(body string) string {
	if prBodyCloseRegex != nil {
		matchString := prBodyCloseRegex.FindString(body)
		return matchString
	}
	return ""
}
