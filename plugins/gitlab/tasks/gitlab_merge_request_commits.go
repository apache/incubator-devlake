package tasks

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"gorm.io/gorm/clause"
)

type ApiMergeRequestCommitResponse []GitlabMergeRequestCommit

type GitlabMergeRequestCommit struct {
	CommitId       string `json:"id"`
	Title          string
	Message        string
	ProjectId      int
	ShortId        string           `json:"short_id"`
	AuthorName     string           `json:"author_name"`
	AuthorEmail    string           `json:"author_email"`
	AuthoredDate   core.Iso8601Time `json:"authored_date"`
	CommitterName  string           `json:"committer_name"`
	CommitterEmail string           `json:"committer_email"`
	CommittedDate  core.Iso8601Time `json:"committed_date"`
	WebUrl         string           `json:"web_url"`
}

var mergeRequestCommitSlice = []models.GitlabMergeRequestCommit{}
var mergeRequestsAssociationSlice = []models.GitlabMergeRequestCommitMergeRequest{}

func CollectMergeRequestCommits(projectId int, mr *models.GitlabMergeRequest) error {
	gitlabApiClient := CreateApiClient()

	getUrl := fmt.Sprintf("projects/%v/merge_requests/%v/commits", projectId, mr.Iid)
	return gitlabApiClient.FetchWithPagination(getUrl, nil, 100,
		func(res *http.Response) error {
			gitlabApiResponse := &ApiMergeRequestCommitResponse{}
			err := core.UnmarshalResponse(res, gitlabApiResponse)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}
			for _, commit := range *gitlabApiResponse {
				gitlabMergeRequestCommit, err := convertMergeRequestCommit(&commit)
				if err != nil {
					return err
				}

				mergeRequestCommitSlice = append(mergeRequestCommitSlice, *gitlabMergeRequestCommit)

				GitlabMergeRequestCommitMergeRequest := &models.GitlabMergeRequestCommitMergeRequest{
					MergeRequestCommitId: commit.CommitId,
					MergeRequestId:       mr.GitlabId,
				}

				mergeRequestsAssociationSlice = append(mergeRequestsAssociationSlice, *GitlabMergeRequestCommitMergeRequest)
			}
			err = saveMergeRequestsCommitsInBatches()
			if err != nil {
				logger.Error("Error: ", err)
				return err
			}
			return nil
		})
}

func saveMergeRequestsCommitsInBatches() error {
	err := lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&mergeRequestCommitSlice).Error
	if err != nil {
		return err
	}
	err = lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&mergeRequestsAssociationSlice).Error
	if err != nil {
		return err
	}
	return nil
}

func convertMergeRequestCommit(commit *GitlabMergeRequestCommit) (*models.GitlabMergeRequestCommit, error) {
	gitlabMergeRequestCommit := &models.GitlabMergeRequestCommit{
		CommitId:       commit.CommitId,
		Title:          commit.Title,
		Message:        commit.Message,
		ShortId:        commit.ShortId,
		AuthorName:     commit.AuthorName,
		AuthorEmail:    commit.AuthorEmail,
		AuthoredDate:   commit.AuthoredDate.ToTime(),
		CommitterName:  commit.CommitterName,
		CommitterEmail: commit.CommitterEmail,
		CommittedDate:  commit.CommittedDate.ToTime(),
		WebUrl:         commit.WebUrl,
	}
	return gitlabMergeRequestCommit, nil
}
