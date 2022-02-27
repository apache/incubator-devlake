package tasks

import (
	"context"
	"fmt"
	"github.com/merico-dev/lake/errors"
	lakeModels "github.com/merico-dev/lake/models"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
	"regexp"
	"strconv"
	"strings"
)

var prBodyCloseRegex *regexp.Regexp
var prBodyClose string

func init() {
	//can not assign config.V.GetString("GITHUB_PR_BODY_CLOSE")
	prBodyClose = ""
}

func EnrichGithubPullRequestIssue(ctx context.Context, repoId int, owner string, repo string) (err error) {
	//only use when we  want to debug ref_bug_stats
	//f, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0x666)
	//defer f.Close()
	//if err != nil {
	//	return nil
	//}
	if len(prBodyClose) > 0 {
		prBodyCloseRegex = regexp.MustCompile(prBodyClose)
	} else {
		patternStr := fmt.Sprintf(
			`(?mi)^.*(fix|close|resolve|fixes|closes|resolves|fixed|closed|resolved)[\s]*(%s\/%s|issue)?(\s)?(((and ){0,1}(#|https:\/\/github\.com\/%s\/%s\/issues\/)\d+[ ]?)+)`,
			owner, repo, owner, repo)
		prBodyCloseRegex = regexp.MustCompile(patternStr)
	}
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
			return errors.TaskCanceled
		default:
		}
		err = lakeModels.Db.ScanRows(cursor, githubPullRequst)
		if err != nil {
			return err
		}

		issueNumberListStr := getCloseIssueId(githubPullRequst.Body)
		//record all matching string, for debug
		//if issueNumberListStr != "" {
		//	_, err := f.WriteString(fmt.Sprintf("%s\n", issueNumberListStr))
		//	if err != nil {
		//		return err
		//	}
		//}

		if issueNumberListStr == "" {
			continue
		}
		//replace https:// to #, then we can deal with later
		if strings.Contains(issueNumberListStr, "https") {
			httpsPrefixPattern := regexp.MustCompile(`https:\/\/github\.com\/\w+\/\w+issues\/`)
			issueNumberListStr = httpsPrefixPattern.ReplaceAllString(issueNumberListStr, "#")

		}
		//split the string by '#'
		firstFilterList := strings.Split(issueNumberListStr, "#")
		issueNumberList := make([]string, len(firstFilterList))

		for _, v := range firstFilterList {
			//split again by ' '
			tmp := strings.Split(v, " ")
			issueNumberList = append(issueNumberList, tmp...)
		}
		for _, issueNumberStr := range issueNumberList {
			issue := &githubModels.GithubIssue{}

			issueNumberStr = strings.TrimSpace(issueNumberStr)
			issueNumber, numFormatErr := strconv.Atoi(issueNumberStr)
			if numFormatErr != nil {
				continue
			}
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
		//if contains not, it may be 'not close/fix/resolve'
		if strings.Contains(matchString, "not") {
			return ""
		}
		return matchString
	}
	return ""
}
