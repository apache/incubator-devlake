package tasks

import (
	"fmt"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"gorm.io/gorm/clause"
)

type ApiMergeRequestNoteResponse []struct {
	GitlabId        int    `json:"id"`
	NoteableId      int    `json:"noteable_id"`
	NoteableIid     int    `json:"noteable_iid"`
	NoteableType    string `json:"noteable_type"`
	Body            string
	GitlabCreatedAt string `json:"created_at"`
	Confidential    bool
	Author          struct {
		Username string `json:"username"`
	}
}

func CollectMergeRequestNotes(projectId int, mrResponse *ApiMergeRequestResponse) error {
	gitlabApiClient := CreateApiClient()

	for _, mr := range *mrResponse {
		getUrl := fmt.Sprintf("projects/%v/merge_requests/%v/notes?system=false", projectId, mr.Iid)
		logger.Info("get URL: ", getUrl)
		res, err := gitlabApiClient.Get(getUrl, nil, nil)
		if err != nil {
			return err
		}

		gitlabApiResponse := &ApiMergeRequestNoteResponse{}

		logger.Info("res", res)

		err = core.UnmarshalResponse(res, gitlabApiResponse)

		if err != nil {
			logger.Error("Error: ", err)
			return nil
		}

		for _, mrNote := range *gitlabApiResponse {
			gitlabMergeRequestNote := &models.GitlabMergeRequestNote{
				GitlabId:        mrNote.GitlabId,
				NoteableId:      mrNote.NoteableId,
				NoteableIid:     mrNote.NoteableIid,
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
	}
	return nil
}
