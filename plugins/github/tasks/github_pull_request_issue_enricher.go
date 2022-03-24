package tasks

import (
	"fmt"
	"github.com/merico-dev/lake/plugins/core"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var EnrichPullRequestIssuesMeta = core.SubTaskMeta{
	Name:             "enrichPullRequestIssues",
	EntryPoint:       EnrichPullRequestIssues,
	EnabledByDefault: true,
	Description:      "Create tool layer table github_pull_request_issues from github_pull_reqeusts",
}

func EnrichPullRequestIssues(taskCtx core.SubTaskContext) (err error) {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	var prBodyCloseRegex *regexp.Regexp
	prBodyClosePattern := taskCtx.GetConfig("GITHUB_PR_BODY_CLOSE_PATTERN")
	numberPrefix := taskCtx.GetConfig("GITHUB_PR_BODY_NUMBER_PREFIX")
	//the pattern before the issue number, sometimes, the issue number is #1098, sometimes it is https://xxx/#1098
	numberPattern := fmt.Sprintf(numberPrefix+`\d+[ ]*)+)`, data.Options.Owner, data.Options.Repo)
	if len(prBodyClosePattern) > 0 {
		prPattern := prBodyClosePattern + numberPattern
		prBodyCloseRegex = regexp.MustCompile(prPattern)
	}
	numberPrefixRegex := regexp.MustCompile(numberPrefix)
	charPattern := regexp.MustCompile(`[a-zA-Z\s,]+`)
	cursor, err := db.Model(&githubModels.GithubPullRequest{}).
		Where("repo_id = ?", repoId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	// iterate all rows

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(githubModels.GithubPullRequest{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			Table: RAW_PULL_REQUEST_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			githubPullRequst := inputRow.(*githubModels.GithubPullRequest)
			results := make([]interface{}, 0, 1)

			//find the matched string in body
			issueNumberListStr := ""

			if prBodyCloseRegex != nil {
				issueNumberListStr = prBodyCloseRegex.FindString(githubPullRequst.Body)
			}

			if issueNumberListStr == "" {
				return nil, nil
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
				err = db.Where("number = ? and repo_id = ?", issueNumber, repoId).
					Limit(1).Find(issue).Error
				if err != nil {
					return nil, err
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
				results = append(results, githubPullRequstIssue)
			}
			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
