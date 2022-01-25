package tasks

import (
	"github.com/merico-dev/lake/config"
	lakeModels "github.com/merico-dev/lake/models"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"regexp"
)

var prType = config.V.GetString("GITHUB_PR_TYPE")
var prComponent = config.V.GetString("GITHUB_PR_COMPONENT")

func EnrichGithubPullRequestWithLabel() (err error) {
	githubPullRequst := &githubModels.GithubPullRequest{}
	cursor, err := lakeModels.Db.Model(&githubPullRequst).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, githubPullRequst)
		githubPullRequst.Type = ""
		githubPullRequst.Component = ""
		if err != nil {
			return err
		}
		var pullRequestLabels []string

		lakeModels.Db.Table("github_issue_labels").
			Where("issue_id = ?", githubPullRequst.GithubId).Order("updated_at ASC").
			Pluck("`label_name`", &pullRequestLabels)

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
	var labelRegex *regexp.Regexp
	// if pr.Type has not been set and prType is set in .env, process the below
	if prType != "" && pr.Type == "" {
		labelRegex = regexp.MustCompile(prType)
	}
	if labelRegex != nil {
		groups := labelRegex.FindStringSubmatch(label)
		if len(groups) > 0 {
			pr.Type = groups[1]
			return
		}
	}

	// if pr.Component has not been set and prComponent is set in .env, process
	if prComponent != "" && pr.Component == "" {
		labelRegex = regexp.MustCompile(prComponent)
	}
	if labelRegex != nil {
		groups := labelRegex.FindStringSubmatch(label)
		if len(groups) > 0 {
			pr.Component = groups[1]
			return
		}
	}
}
