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

type ApiCommitResponse []struct {
	GitlabId       string `json:"id"`
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

func CollectCommits(projectId int, scheduler *utils.WorkerScheduler) error {
	gitlabApiClient := CreateApiClient()

	return gitlabApiClient.FetchWithPaginationAnts(scheduler, fmt.Sprintf("projects/%v/repository/commits?with_stats=true", projectId), 100,
		func(res *http.Response) error {

			gitlabApiResponse := &ApiCommitResponse{}
			err := core.UnmarshalResponse(res, gitlabApiResponse)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}

			for _, value := range *gitlabApiResponse {
				gitlabCommit := &models.GitlabCommit{
					GitlabId:       value.GitlabId,
					Title:          value.Title,
					Message:        value.Message,
					ProjectId:      projectId,
					ShortId:        value.ShortId,
					AuthorName:     value.AuthorName,
					AuthorEmail:    value.AuthorEmail,
					AuthoredDate:   utils.ConvertStringToTime(value.AuthoredDate),
					CommitterName:  value.CommitterName,
					CommitterEmail: value.CommitterEmail,
					CommittedDate:  utils.ConvertStringToTime(value.CommittedDate),
					WebUrl:         value.WebUrl,
					Additions:      value.Stats.Additions,
					Deletions:      value.Stats.Deletions,
					Total:          value.Stats.Total,
				}

				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&gitlabCommit).Error

				if err != nil {
					logger.Error("Could not upsert: ", err)
				}
			}

			return nil
		})
}
