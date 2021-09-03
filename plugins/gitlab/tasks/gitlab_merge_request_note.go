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

type ApiMergeRequestNoteResponse []struct {
	GitlabId        int    `json:"id"`
	NoteableId      int    `json:"noteable_id"`
	MergeRequestIid int    `json:"noteable_iid"`
	NoteableType    string `json:"noteable_type"`
	Body            string
	GitlabCreatedAt string `json:"created_at"`
	Confidential    bool
	Author          struct {
		Username string `json:"username"`
	}
}

func CollectMergeRequestNotes(projectId int, mrId int) error {
	gitlabApiClient := CreateApiClient()

	getUrl := fmt.Sprintf("projects/%v/merge_requests/%v/notes?system=false", projectId, mrId)
	return gitlabApiClient.FetchWithPagination(getUrl, "100",
		func(res *http.Response) error {

			gitlabApiResponse := &ApiMergeRequestNoteResponse{}
			err := core.UnmarshalResponse(res, gitlabApiResponse)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}

			for _, mrNote := range *gitlabApiResponse {
				gitlabMergeRequestNote := &models.GitlabMergeRequestNote{
					GitlabId:        mrNote.GitlabId,
					NoteableId:      mrNote.NoteableId,
					MergeRequestId:  mrNote.MergeRequestIid,
					NoteableType:    mrNote.NoteableType,
					AuthorUsername:  mrNote.Author.Username,
					Body:            mrNote.Body,
					GitlabCreatedAt: mrNote.GitlabCreatedAt,
					Confidential:    mrNote.Confidential,
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
}
