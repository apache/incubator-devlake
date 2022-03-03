package tasks

import (
	"context"
	"fmt"
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/utils"
	"gorm.io/gorm/clause"
	"net/http"
)

type MergeRequestNote struct {
	GitlabId        int    `json:"id"`
	MergeRequestId  int    `json:"noteable_id"`
	MergeRequestIid int    `json:"noteable_iid"`
	NoteableType    string `json:"noteable_type"`
	Body            string
	GitlabCreatedAt core.Iso8601Time `json:"created_at"`
	Confidential    bool
	Resolvable      bool `json:"resolvable"`
	System          bool `json:"system"`
	Author          struct {
		Username string `json:"username"`
	}
}
type ApiMergeRequestNoteResponse []MergeRequestNote

func CollectMergeRequestNotes(ctx context.Context, projectId int, rateLimitPerSecondInt int, gitlabApiClient *GitlabApiClient) error {
	scheduler, err := utils.NewWorkerScheduler(rateLimitPerSecondInt*2, rateLimitPerSecondInt, ctx)
	if err != nil {
		return nil
	}
	defer scheduler.Release()
	gitlabMr := &models.GitlabMergeRequest{}
	cursor, err := lakeModels.Db.Model(gitlabMr).Where("project_id = ?", projectId).Rows()
	if err != nil {
		return nil
	}
	defer cursor.Close()

	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, gitlabMr)
		if err != nil {
			return nil
		}
		getUrl := fmt.Sprintf("projects/%v/merge_requests/%v/notes?system=false", projectId, gitlabMr.Iid)
		err = scheduler.Submit(func() error {
			return gitlabApiClient.FetchWithPagination(getUrl, nil, 100,
				func(res *http.Response) error {

					gitlabApiResponse := &ApiMergeRequestNoteResponse{}
					err = core.UnmarshalResponse(res, gitlabApiResponse)

					if err != nil {
						logger.Error("Error: ", err)
						return nil
					}

					for _, mrNote := range *gitlabApiResponse {
						gitlabMergeRequestNote, err := convertMergeRequestNote(&mrNote)
						if err != nil {
							return err
						}
						err = lakeModels.Db.Clauses(clause.OnConflict{
							UpdateAll: true,
						}).Create(&gitlabMergeRequestNote).Error

						if err != nil {
							logger.Error("Could not upsert: ", err)
							return err
						}
					}
					return nil
				})
		})
		if err != nil {
			return err
		}
	}
	scheduler.WaitUntilFinish()
	return nil
}
func convertMergeRequestNote(mrNote *MergeRequestNote) (*models.GitlabMergeRequestNote, error) {
	gitlabMergeRequestNote := &models.GitlabMergeRequestNote{
		GitlabId:        mrNote.GitlabId,
		MergeRequestId:  mrNote.MergeRequestId,
		MergeRequestIid: mrNote.MergeRequestIid,
		NoteableType:    mrNote.NoteableType,
		AuthorUsername:  mrNote.Author.Username,
		Body:            mrNote.Body,
		GitlabCreatedAt: mrNote.GitlabCreatedAt.ToTime(),
		Confidential:    mrNote.Confidential,
		Resolvable:      mrNote.Resolvable,
		System:          mrNote.System,
	}
	return gitlabMergeRequestNote, nil
}
