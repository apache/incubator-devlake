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

	var commitsToSave = ConcurrentCommits{}
	var projectCommitsToSave = ConcurrentProjectCommits{}
	var usersToSave = ConcurrentUsers{}

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

				commitsToSave.Append(gitlabCommit)
				// create project/commits relationship
				gitlabProjectCommit.CommitSha = gitlabCommit.Sha
				projectCommitsToSave.Append(gitlabProjectCommit)

				addUsersToSlice(*gitlabCommit, &usersToSave)
			}
			return nil
		})
	// listen for the last ants submission before saving the data
	<-finish
	// when we receive the message, we have to wait for the scheduler to finish its
	// tasks before we save data
	fmt.Println("INFO >>> all done collecting!")
	scheduler.WaitUntilFinish()
	fmt.Println("INFO >>> saving gitlab_commits")
	err := saveSlice(commitsToSave.commits)
	if err != nil {
		logger.Error("Error: ", err)
		return err
	}
	fmt.Println("INFO >>> saving gitlab_project_commits")
	err = saveSlice(projectCommitsToSave.projectCommits)
	if err != nil {
		logger.Error("Error: ", err)
		return err
	}
	fmt.Println("INFO >>> saving gitlab_users")
	err = saveSlice(usersToSave.users)
	if err != nil {
		logger.Error("Error: ", err)
		return err
	}
	return nil
}

func saveSlice(data interface{}) error {
	err := lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func addUsersToSlice(commit models.GitlabCommit, usersToSave *ConcurrentUsers) {
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
	for _, user := range usersToSave.users {
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
		usersToSave.Append(gitlabAuthor)
	}

	if committerIsDifferent && !committerExists {
		usersToSave.Append(gitlabCommitter)
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
