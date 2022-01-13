package tasks

import (
	"fmt"
	"net/http"
	"time"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"gorm.io/gorm/clause"
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

func FindEarliestNote(notes *ApiMergeRequestNoteResponse) (*MergeRequestNote, error) {
	var earliestNote *MergeRequestNote

	earliestTime := time.Now()
	for _, note := range *notes {
		if note.System || !note.Resolvable {
			continue
		}
		noteTime := note.GitlabCreatedAt.ToTime()
		if noteTime.Before(earliestTime) {
			earliestTime = noteTime
			earliestNote = &note
		}
	}
	return earliestNote, nil
}

// we need a metric that measures a merge request duration as the time from first comment to MR close
func updateMergeRequestWithFirstCommentTime(notes *ApiMergeRequestNoteResponse, mr *models.GitlabMergeRequest) error {
	earliestNote, err := FindEarliestNote(notes)
	if err != nil {
		return err
	}
	if earliestNote != nil {
		t := earliestNote.GitlabCreatedAt.ToTime()
		mr.FirstCommentTime = &t

		err = lakeModels.Db.Model(&mr).Where("gitlab_id = ?", mr.GitlabId).Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Update("first_comment_time", mr.FirstCommentTime).Error

		if err != nil {
			logger.Error("Could not upsert: ", err)
			return err
		}
	}
	return nil
}

var mergeRequestNotesSlice = []models.GitlabMergeRequestNote{}

func CollectMergeRequestNotes(projectId int, mr *models.GitlabMergeRequest) error {
	gitlabApiClient := CreateApiClient()

	getUrl := fmt.Sprintf("projects/%v/merge_requests/%v/notes?system=false", projectId, mr.Iid)
	return gitlabApiClient.FetchWithPagination(getUrl, nil, 100,
		func(res *http.Response) error {

			gitlabApiResponse := &ApiMergeRequestNoteResponse{}
			err := core.UnmarshalResponse(res, gitlabApiResponse)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}

			for _, mrNote := range *gitlabApiResponse {
				gitlabMergeRequestNote, err := convertMergeRequestNote(&mrNote)
				if err != nil {
					return err
				}
				mergeRequestNotesSlice = append(mergeRequestNotesSlice, *gitlabMergeRequestNote)
			}

			mergeRequestUpdateErr := updateMergeRequestWithFirstCommentTime(gitlabApiResponse, mr)
			if mergeRequestUpdateErr != nil {
				return err
			}
			err = saveMergeRequestsNotesInBatches()
			if err != nil {
				logger.Error("Error: ", err)
				return err
			}
			return nil
		})
}

func saveMergeRequestsNotesInBatches() error {
	err := lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&mergeRequestNotesSlice).Error
	if err != nil {
		return err
	}
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
