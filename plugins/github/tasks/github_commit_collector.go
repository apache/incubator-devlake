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

type ApiSingleCommitResponse struct {
	Stats struct {
		Additions int
		Deletions int
	}
}

func CollectCommits(owner string, repo string, repoId int, apiClient *GithubApiClient) error {
	getUrl := fmt.Sprintf("repos/%v/%v/commits", owner, repo)
	return apiClient.FetchPages(getUrl, nil, 100,
		func(res *http.Response) error {
			githubApiResponse := &ApiCommitsResponse{}
			err := core.UnmarshalResponse(res, githubApiResponse)
			if err != nil || res.StatusCode == 401 {
				return err
			}
			repoCommit := &models.GithubRepoCommit{RepoId: repoId}
			for _, commit := range *githubApiResponse {
				githubCommit, err := convertGithubCommit(&commit)
				if err != nil {
					return err
				}
				// save author and committer
				if commit.Author != nil {
					githubCommit.AuthorId = commit.Author.Id
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&commit.Author).Error
					if err != nil {
						return err
					}
				}
				if commit.Committer != nil {
					githubCommit.CommitterId = commit.Committer.Id
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&commit.Committer).Error
					if err != nil {
						return err
					}
				}
				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&githubCommit).Error
				if err != nil {
					return err
				}
				// save repo / commit relationship
				repoCommit.CommitSha = commit.Sha
				err = lakeModels.Db.Clauses(clause.OnConflict{
					DoNothing: true,
				}).Create(repoCommit).Error
				if err != nil {
					return err
				}
			}
			return nil
		})
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

func CollectCommitsStat(
	owner string,
	repo string,
	repoId int,
	scheduler *utils.WorkerScheduler,
	apiClient *GithubApiClient,
) error {
	cursor, err := lakeModels.Db.Table("github_commits gc").
		Joins(`left join github_repo_commits grc on (
			grc.commit_sha = gc.sha
		)`).
		Select("gc.*").
		Where("grc.repo_id = ?", repoId).
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	// TODO: this still loading all rows into memory, to be optimized
	for cursor.Next() {
		commit := &models.GithubCommit{}
		err = lakeModels.Db.ScanRows(cursor, commit)
		if err != nil {
			return err
		}
		err = scheduler.Submit(func() error {
			// This call is to update the details of the individual pull request with additions / deletions / etc.
			err := CollectCommit(owner, repo, repoId, commit, apiClient)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil
		}
	}

	scheduler.WaitUntilFinish()
	return nil
}

// for addtions and deletions
func CollectCommit(
	owner string,
	repo string,
	repoId int,
	commit *models.GithubCommit,
	apiClient *GithubApiClient,
) error {
	getUrl := fmt.Sprintf("repos/%v/%v/commits/%v", owner, repo, commit.Sha)
	res, getErr := apiClient.Get(getUrl, nil, nil)
	if getErr != nil {
		logger.Error("GET Error: ", getErr)
		return getErr
	}

	githubApiResponse := &ApiSingleCommitResponse{}
	unmarshalErr := core.UnmarshalResponse(res, githubApiResponse)
	if unmarshalErr != nil {
		logger.Error("Error: ", unmarshalErr)
		return unmarshalErr
	}
	dbErr := lakeModels.Db.Model(&commit).Updates(models.GithubCommit{
		Additions: githubApiResponse.Stats.Additions,
		Deletions: githubApiResponse.Stats.Deletions,
	}).Error
	if dbErr != nil {
		logger.Error("Could not update: ", dbErr)
	}
	return nil
}
