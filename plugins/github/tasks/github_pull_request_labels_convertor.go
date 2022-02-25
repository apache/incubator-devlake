package tasks

import (
	"context"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func ConvertPullRequestLabels(ctx context.Context, repoId int) error {
	githubPullRequestLabel := &githubModels.GithubPullRequestLabel{}
	cursor, err := lakeModels.Db.Model(githubPullRequestLabel).
		Joins(`left join github_pull_requests on github_pull_requests.github_id = github_pull_request_labels.pull_id`).
		Where("github_pull_requests.repo_id = ?", repoId).
		Order("pull_id ASC").
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	domainIdGeneratorPr := didgen.NewDomainIdGenerator(&githubModels.GithubPullRequest{})
	lastPrId := 0
	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, githubPullRequestLabel)
		if err != nil {
			return err
		}
		pullRequestId := domainIdGeneratorPr.Generate(githubPullRequestLabel.PullId)
		if lastPrId != githubPullRequestLabel.PullId {
			// Clean up old data
			err := lakeModels.Db.Where("pull_request_id = ?",
				pullRequestId).Delete(&code.PullRequestLabel{}).Error
			if err != nil {
				return err
			}
			lastPrId = githubPullRequestLabel.PullId
		}

		pullRequestLabel := &code.PullRequestLabel{
			PullRequestId: pullRequestId,
			LabelName:     githubPullRequestLabel.LabelName,
		}
		err = lakeModels.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(pullRequestLabel).Error
		if err != nil {
			return err
		}
	}
	return nil
}
