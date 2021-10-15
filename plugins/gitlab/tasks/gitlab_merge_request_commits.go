package tasks

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/utils"
	"gorm.io/gorm/clause"
)

type ApiMergeRequestCommitResponse []struct {
	CommitId       string `json:"id"`
	Title          string
	Message        string
	ProjectId      int
	ShortId        string `json:"short_id"`
	AuthorName     string `json:"author_name"`
	AuthorEmail    string `json:"author_email"`
	AuthoredDate   string `json:"authored_date"`
	CommitterName  string `json:"committer_name"`
	CommitterEmail string `json:"committer_email"`
	CommittedDate  string `json:"committed_date"`
	WebUrl         string `json:"web_url"`
	Stats          struct {
		Additions int
		Deletions int
		Total     int
	}
}

func CollectMergeRequestCommits(projectId int, mr *models.GitlabMergeRequest) error {
	gitlabApiClient := CreateApiClient()

	getUrl := fmt.Sprintf("projects/%v/merge_requests/%v/commits", projectId, mr.Iid)
	return gitlabApiClient.FetchWithPagination(getUrl, 100,
		func(res *http.Response) error {
			gitlabApiResponse := &ApiMergeRequestCommitResponse{}
			err := core.UnmarshalResponse(res, gitlabApiResponse)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}
			for _, commit := range *gitlabApiResponse {
				gitlabMergeRequestCommit := &models.GitlabMergeRequestCommit{
					CommitId:       commit.CommitId,
					Title:          commit.Title,
					Message:        commit.Message,
					ShortId:        commit.ShortId,
					AuthorName:     commit.AuthorName,
					AuthorEmail:    commit.AuthorEmail,
					AuthoredDate:   utils.ConvertStringToSqlNullTime(commit.AuthoredDate),
					CommitterName:  commit.CommitterName,
					CommitterEmail: commit.CommitterEmail,
					CommittedDate:  utils.ConvertStringToSqlNullTime(commit.CommittedDate),
					WebUrl:         commit.WebUrl,
					Additions:      commit.Stats.Additions,
					Deletions:      commit.Stats.Deletions,
					Total:          commit.Stats.Total,
				}
				result := lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&gitlabMergeRequestCommit)

				if result.Error != nil {
					logger.Error("Could not upsert: ", result.Error)
				}
				GitlabMergeRequestCommitMergeRequest := &models.GitlabMergeRequestCommitMergeRequest{
					MergeRequestCommitId: commit.CommitId,
					MergeRequestId:       mr.GitlabId,
				}
				result = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&GitlabMergeRequestCommitMergeRequest)

				if result.Error != nil {
					logger.Error("Could not upsert: ", result.Error)
				}
			}

			return nil
		})
}
