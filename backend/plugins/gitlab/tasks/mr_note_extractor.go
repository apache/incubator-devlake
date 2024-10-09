/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tasks

import (
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiMrNotesMeta)
}

type MergeRequestNote struct {
	GitlabId        int    `json:"id"`
	MergeRequestId  int    `json:"noteable_id"`
	MergeRequestIid int    `json:"noteable_iid"`
	NoteableType    string `json:"noteable_type"`
	Body            string
	GitlabCreatedAt common.Iso8601Time `json:"created_at"`
	Confidential    bool
	Resolvable      bool `json:"resolvable"`
	System          bool `json:"system"`
	Author          struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
	}
	Type string `json:"type"`
}

var ExtractApiMrNotesMeta = plugin.SubTaskMeta{
	Name:             "Extract MR Notes",
	EntryPoint:       ExtractApiMergeRequestsNotes,
	EnabledByDefault: true,
	Description:      "Extract raw merge requests notes data into tool layer table GitlabMrNote",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
	Dependencies:     []*plugin.SubTaskMeta{&CollectApiMrNotesMeta},
}

func ExtractApiMergeRequestsNotes(subtaskCtx plugin.SubTaskContext) errors.Error {
	subtaskCommonArgs, data := CreateSubtaskCommonArgs(subtaskCtx, RAW_MERGE_REQUEST_NOTES_TABLE)

	extractor, err := api.NewStatefulApiExtractor(&api.StatefulApiExtractorArgs[MergeRequestNote]{
		SubtaskCommonArgs: subtaskCommonArgs,
		Extract: func(mrNote *MergeRequestNote, row *api.RawData) ([]interface{}, errors.Error) {
			toolMrNote, err := convertMergeRequestNote(mrNote)
			toolMrNote.ConnectionId = data.Options.ConnectionId
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, 2)
			if !toolMrNote.IsSystem || toolMrNote.Body == "approved this merge request" || toolMrNote.Body == "unapproved this merge request" {
				toolMrComment := &models.GitlabMrComment{
					GitlabId:        toolMrNote.GitlabId,
					MergeRequestId:  toolMrNote.MergeRequestId,
					MergeRequestIid: toolMrNote.MergeRequestIid,
					Body:            toolMrNote.Body,
					AuthorUserId:    toolMrNote.AuthorUserId,
					AuthorUsername:  toolMrNote.AuthorUsername,
					GitlabCreatedAt: toolMrNote.GitlabCreatedAt,
					Resolvable:      toolMrNote.Resolvable,
					Type:            toolMrNote.Type,
					ConnectionId:    data.Options.ConnectionId,
				}
				if toolMrNote.Body == "approved this merge request" {
					toolMrComment.Type = "REVIEW"
				}
				results = append(results, toolMrComment)
			}
			results = append(results, toolMrNote)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertMergeRequestNote(mrNote *MergeRequestNote) (*models.GitlabMrNote, errors.Error) {
	GitlabMrNote := &models.GitlabMrNote{
		GitlabId:        mrNote.GitlabId,
		AuthorUserId:    mrNote.Author.Id,
		MergeRequestId:  mrNote.MergeRequestId,
		MergeRequestIid: mrNote.MergeRequestIid,
		NoteableType:    mrNote.NoteableType,
		AuthorUsername:  mrNote.Author.Username,
		Body:            mrNote.Body,
		GitlabCreatedAt: mrNote.GitlabCreatedAt.ToTime(),
		Confidential:    mrNote.Confidential,
		Resolvable:      mrNote.Resolvable,
		IsSystem:        mrNote.System,
		Type:            mrNote.Type,
	}
	return GitlabMrNote, nil
}
