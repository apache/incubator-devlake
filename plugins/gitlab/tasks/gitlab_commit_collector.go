package tasks

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/utils"
	"gorm.io/gorm/clause"
)

type ApiCommitResponse []GitlabApiCommit

type GitlabApiCommit struct {
	GitlabId       string `json:"id"`
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
	Stats          struct {
		Additions int
		Deletions int
		Total     int
	}
}

func CollectCommits(projectId int, scheduler *utils.WorkerScheduler) error {
	// Temporary storage before DB save
	var commitsToSave = []models.GitlabCommit{}
	var projectCommitsToSave = []models.GitlabProjectCommit{}
	var usersToSave = []models.GitlabUser{}

	gitlabApiClient := CreateApiClient()
	relativePath := fmt.Sprintf("projects/%v/repository/commits", projectId)
	queryParams := &url.Values{}
	queryParams.Set("with_stats", "true")
	finish := make(chan bool) // This will tell us when there are no more pages to fetch
	go gitlabApiClient.FetchWithPaginationAnts(finish, scheduler, relativePath, queryParams, 100,
		func(res *http.Response) error { // handles the response (with 100 results) from API

			gitlabApiResponse := &ApiCommitResponse{}
			err := core.UnmarshalResponse(res, gitlabApiResponse)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}
			gitlabProjectCommit := &models.GitlabProjectCommit{GitlabProjectId: projectId}
			for _, gitlabApiCommit := range *gitlabApiResponse {
				gitlabCommit, err := convertCommit(&gitlabApiCommit, projectId)
				if err != nil {
					return err
				}

				commitsToSave = append(commitsToSave, *gitlabCommit)

				// create project/commits relationship
				gitlabProjectCommit.CommitSha = gitlabCommit.Sha
				projectCommitsToSave = append(projectCommitsToSave, *gitlabProjectCommit)

				addUsersToSlice(*gitlabCommit, &usersToSave)
			}
			return nil
		})
	// listen for the last ants submission before saving the data
	<-finish
	// when we receive the message, we have to wait for the scheduler to finish its
	// tasks before we save data
	scheduler.WaitUntilFinish()
	err := saveSlice("gitlab_commits", commitsToSave)
	if err != nil {
		logger.Error("Error: ", err)
		return err
	}
	err = saveSlice("gitlab_project_commits", projectCommitsToSave)
	if err != nil {
		logger.Error("Error: ", err)
		return err
	}
	err = saveSlice("gitlab_users", usersToSave)
	if err != nil {
		logger.Error("Error: ", err)
		return err
	}
	return nil
}

func saveSlice(table string, data interface{}) error {
	fmt.Println("KEVIN >>> saving data...", table)
	err := lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func addUsersToSlice(commit models.GitlabCommit, usersToSave *[]models.GitlabUser) {
	authorExists := false
	committerExists := false
	gitlabAuthor := &models.GitlabUser{
		Email: commit.AuthorEmail,
		Name:  commit.AuthorName,
	}
	gitlabCommitter := &models.GitlabUser{
		Email: commit.CommitterEmail,
		Name:  commit.CommitterName,
	}
	committerIsDifferent := false

	if commit.AuthorEmail != commit.CommitterEmail {
		committerIsDifferent = true
	}
	for _, user := range *usersToSave {
		if user.Email == gitlabAuthor.Email {
			authorExists = true
		}
		if committerIsDifferent {
			if user.Email == gitlabCommitter.Email {
				committerExists = true
			}
		}
	}
	if !authorExists {
		*usersToSave = append(*usersToSave, *gitlabAuthor)
	}

	if committerIsDifferent && !committerExists {
		*usersToSave = append(*usersToSave, *gitlabCommitter)
	}
}

// Convert the API response to our DB model instance
func convertCommit(commit *GitlabApiCommit, projectId int) (*models.GitlabCommit, error) {
	gitlabCommit := &models.GitlabCommit{
		Sha:            commit.GitlabId,
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
		Additions:      commit.Stats.Additions,
		Deletions:      commit.Stats.Deletions,
		Total:          commit.Stats.Total,
	}
	return gitlabCommit, nil
}
