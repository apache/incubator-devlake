package tasks

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/utils"
	"gorm.io/gorm/clause"
)

type ApiCommitsResponse []CommitsResponse
type CommitsResponse struct {
	Sha       string `json:"sha"`
	Commit    Commit
	Url       string
	Author    *models.GithubUser
	Committer *models.GithubUser
}

type Commit struct {
	Author struct {
		Name  string
		Email string
		Date  core.Iso8601Time
	}
	Committer struct {
		Name  string
		Email string
		Date  core.Iso8601Time
	}
	Message string
}

// Store the data in slices so we can batch insert later
var commitSlice = []models.GithubCommit{}
var repoCommitsSlice = []models.GithubRepoCommit{}
var usersSlice = []models.GithubUser{}

func CollectCommits(owner string, repositoryName string, repositoryId int, scheduler *utils.WorkerScheduler, githubApiClient *GithubApiClient) error {
	getUrl := fmt.Sprintf("repos/%v/%v/commits", owner, repositoryName)
	githubApiClient.FetchWithPaginationAnts(getUrl, nil, 100, 20, scheduler,
		func(res *http.Response) error {
			githubApiResponse := &ApiCommitsResponse{}
			err := core.UnmarshalResponse(res, githubApiResponse)
			if err != nil || res.StatusCode == 401 {
				return err
			}
			repoCommit := &models.GithubRepoCommit{GithubRepoId: repositoryId}
			repoCommitsSlice = append(repoCommitsSlice, *repoCommit)
			fmt.Println("KEVIN >>> len(githubApiResponse): ", len(*githubApiResponse))
			for i, commit := range *githubApiResponse {
				fmt.Println("KEVIN >>> i", i)
				githubCommit, err := convertGithubCommit(&commit)
				if err != nil {
					return err
				}

				commitSlice = append(commitSlice, *githubCommit)

				repoCommit.CommitSha = commit.Sha
				repoCommitsSlice = append(repoCommitsSlice, *repoCommit)

				if commit.Author != nil {
					usersSlice = append(usersSlice, *commit.Author)
				}
				if commit.Committer != nil {
					usersSlice = append(usersSlice, *commit.Committer)
				}
			}
			return nil
		})

	err := insertDataInBatches()
	if err != nil {
		logger.Error("Could not upsert: ", err)
		return err
	} else {
		return nil
	}
}

func insertDataInBatches() error {
	err := lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&commitSlice).Error
	if err != nil {
		return err
	}
	err = lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&repoCommitsSlice).Error
	if err != nil {
		return err
	}
	err = lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&usersSlice).Error
	if err != nil {
		return err
	}
	return nil
}

func convertGithubCommit(commit *CommitsResponse) (*models.GithubCommit, error) {
	githubCommit := &models.GithubCommit{
		Sha:            commit.Sha,
		Message:        commit.Commit.Message,
		AuthorName:     commit.Commit.Author.Name,
		AuthorEmail:    commit.Commit.Author.Email,
		AuthoredDate:   commit.Commit.Author.Date.ToTime(),
		CommitterName:  commit.Commit.Committer.Name,
		CommitterEmail: commit.Commit.Committer.Email,
		CommittedDate:  commit.Commit.Committer.Date.ToTime(),
		Url:            commit.Url,
	}
	return githubCommit, nil
}
