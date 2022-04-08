package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	gitlabModels "github.com/merico-dev/lake/plugins/gitlab/models"
	"gorm.io/gorm/clause"
)

func ConvertMergeRequestCommits(projectId int) error {
	gitlabMergeRequestCommit := &gitlabModels.GitlabMergeRequestCommit{}
	cursor, err := lakeModels.Db.Model(&gitlabMergeRequestCommit).
		Joins(`left join gitlab_merge_requests on gitlab_merge_requests.gitlab_id = gitlab_merge_request_commits.merge_request_id`).
		Where("gitlab_merge_requests.project_id = ?", projectId).
		Order("merge_request_id ASC").Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	var pullRequestId int
	domainPullRequestId := ""
	domainIdGenerator := didgen.NewDomainIdGenerator(&gitlabModels.GitlabMergeRequest{})

	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, gitlabMergeRequestCommit)

		if pullRequestId != gitlabMergeRequestCommit.MergeRequestId {
			domainPullRequestId = domainIdGenerator.Generate(gitlabMergeRequestCommit.MergeRequestId)
			err := lakeModels.Db.Where("pull_request_id = ?",
				domainPullRequestId).Delete(&code.PullRequestCommit{}).Error
			if err != nil {
				return err
			}
			pullRequestId = gitlabMergeRequestCommit.MergeRequestId
		}
		if err != nil {
			return err
		}

		domainPrcommit := &code.PullRequestCommit{
			CommitSha:     gitlabMergeRequestCommit.CommitSha,
			PullRequestId: domainPullRequestId,
		}

		err = lakeModels.Db.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(domainPrcommit).Error
		if err != nil {
			return err
		}
	}
	return nil
}
