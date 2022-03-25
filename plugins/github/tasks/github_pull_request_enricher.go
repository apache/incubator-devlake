package tasks

import (
	"context"
	"github.com/merico-dev/lake/errors"
	"regexp"

	"github.com/merico-dev/lake/config"
	lakeModels "github.com/merico-dev/lake/models"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
)

func EnrichPullRequests(ctx context.Context, repoId int) (err error) {
	var labelTypeRegex *regexp.Regexp
	var labelComponentRegex *regexp.Regexp

	v := config.GetConfig()
	var prType = v.GetString("GITHUB_PR_TYPE")
	var prComponent = v.GetString("GITHUB_PR_COMPONENT")
	if len(prType) > 0 {
		labelTypeRegex = regexp.MustCompile(prType)
	}
	if len(prComponent) > 0 {
		labelComponentRegex = regexp.MustCompile(prComponent)
	}

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
			if labelTypeRegex != nil {
				groups := labelTypeRegex.FindStringSubmatch(pullRequestLabel)
				if len(groups) > 0 {
					githubPullRequst.Type = groups[1]
					continue
				}
			}

			// if pr.Component has not been set and prComponent is set in .env, process
			if labelComponentRegex != nil {
				groups := labelComponentRegex.FindStringSubmatch(pullRequestLabel)
				if len(groups) > 0 {
					githubPullRequst.Component = groups[1]
					continue
				}
			}
		}

		err = lakeModels.Db.Save(githubPullRequst).Error
		if err != nil {
			return err
		}
	}
	return nil
}
