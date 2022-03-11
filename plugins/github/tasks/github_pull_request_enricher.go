package tasks

import (
	"context"
	"regexp"

	"github.com/merico-dev/lake/config"
	lakeModels "github.com/merico-dev/lake/models"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
)

var labelTypeRegex *regexp.Regexp
var labelComponentRegex *regexp.Regexp

func init() {
	V := config.GetConfig()
	var prType = V.GetString("GITHUB_PR_TYPE")
	var prComponent = V.GetString("GITHUB_PR_COMPONENT")
	if len(prType) > 0 {
		labelTypeRegex = regexp.MustCompile(prType)
	}
	if len(prComponent) > 0 {
		labelComponentRegex = regexp.MustCompile(prComponent)
	}
}

func EnrichPullRequests(ctx context.Context, repoId int) (err error) {
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
			return ctx.Err()
		default:
		}
		err = lakeModels.Db.ScanRows(cursor, githubPullRequst)
		if err != nil {
			return err
		}
		githubPullRequst.Type = ""
		githubPullRequst.Component = ""
		var pullRequestLabels []string
		err = lakeModels.Db.Table("github_pull_request_labels").
			Where("pull_id = ?", githubPullRequst.GithubId).
			Pluck("`label_name`", &pullRequestLabels).Error
		if err != nil {
			return err
		}

		for _, pullRequestLabel := range pullRequestLabels {
			setPullRequestLabel(pullRequestLabel, githubPullRequst)
		}

		err = lakeModels.Db.Save(githubPullRequst).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func setPullRequestLabel(label string, pr *githubModels.GithubPullRequest) {
	// if pr.Type has not been set and prType is set in .env, process the below
	if labelTypeRegex != nil {
		groups := labelTypeRegex.FindStringSubmatch(label)
		if len(groups) > 0 {
			pr.Type = groups[1]
			return
		}
	}

	// if pr.Component has not been set and prComponent is set in .env, process
	if labelComponentRegex != nil {
		groups := labelComponentRegex.FindStringSubmatch(label)
		if len(groups) > 0 {
			pr.Component = groups[1]
			return
		}
	}
}
