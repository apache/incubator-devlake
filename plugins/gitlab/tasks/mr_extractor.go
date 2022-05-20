package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/plugins/helper"
)

type MergeRequestRes struct {
	GitlabId        int `json:"id"`
	Iid             int
	ProjectId       int `json:"project_id"`
	SourceProjectId int `json:"source_project_id"`
	TargetProjectId int `json:"target_project_id"`
	State           string
	Title           string
	Description     string
	WebUrl          string              `json:"web_url"`
	UserNotesCount  int                 `json:"user_notes_count"`
	WorkInProgress  bool                `json:"work_in_progress"`
	SourceBranch    string              `json:"source_branch"`
	TargetBranch    string              `json:"target_branch"`
	GitlabCreatedAt helper.Iso8601Time  `json:"created_at"`
	MergedAt        *helper.Iso8601Time `json:"merged_at"`
	ClosedAt        *helper.Iso8601Time `json:"closed_at"`
	MergeCommitSha  string              `json:"merge_commit_sha"`
	MergedBy        struct {
		Username string `json:"username"`
	} `json:"merged_by"`
	Author struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
	}
	Reviewers        []Reviewer
	FirstCommentTime helper.Iso8601Time
}

type Reviewer struct {
	GitlabId       int `json:"id"`
	MergeRequestId int
	Name           string
	Username       string
	State          string
	AvatarUrl      string `json:"avatar_url"`
	WebUrl         string `json:"web_url"`
}

var ExtractApiMergeRequestsMeta = core.SubTaskMeta{
	Name:             "extractApiMergeRequests",
	EntryPoint:       ExtractApiMergeRequests,
	EnabledByDefault: true,
	Description:      "Extract raw merge requests data into tool layer table GitlabMergeRequest and GitlabReviewer",
}

func ExtractApiMergeRequests(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			mr := &MergeRequestRes{}
			err := json.Unmarshal(row.Data, mr)
			if err != nil {
				return nil, err
			}

			gitlabMergeRequest, err := convertMergeRequest(mr)
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, len(mr.Reviewers)+1)

			results = append(results, gitlabMergeRequest)

			for _, reviewer := range mr.Reviewers {
				gitlabReviewer := &models.GitlabReviewer{
					GitlabId:       reviewer.GitlabId,
					MergeRequestId: mr.GitlabId,
					ProjectId:      data.Options.ProjectId,
					Username:       reviewer.Username,
					Name:           reviewer.Name,
					State:          reviewer.State,
					AvatarUrl:      reviewer.AvatarUrl,
					WebUrl:         reviewer.WebUrl,
				}
				results = append(results, gitlabReviewer)
			}

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertMergeRequest(mr *MergeRequestRes) (*models.GitlabMergeRequest, error) {
	gitlabMergeRequest := &models.GitlabMergeRequest{
		GitlabId:         mr.GitlabId,
		Iid:              mr.Iid,
		ProjectId:        mr.ProjectId,
		SourceProjectId:  mr.SourceProjectId,
		TargetProjectId:  mr.TargetProjectId,
		State:            mr.State,
		Title:            mr.Title,
		Description:      mr.Description,
		WebUrl:           mr.WebUrl,
		UserNotesCount:   mr.UserNotesCount,
		WorkInProgress:   mr.WorkInProgress,
		SourceBranch:     mr.SourceBranch,
		TargetBranch:     mr.TargetBranch,
		MergeCommitSha:   mr.MergeCommitSha,
		MergedAt:         helper.Iso8601TimeToTime(mr.MergedAt),
		GitlabCreatedAt:  mr.GitlabCreatedAt.ToTime(),
		ClosedAt:         helper.Iso8601TimeToTime(mr.ClosedAt),
		MergedByUsername: mr.MergedBy.Username,
		AuthorUsername:   mr.Author.Username,
		AuthorUserId:     mr.Author.Id,
	}
	return gitlabMergeRequest, nil
}
