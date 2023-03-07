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
	"encoding/json"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/trello/models"
	"time"
)

var _ plugin.SubTaskEntryPoint = ExtractCard

var ExtractCardMeta = plugin.SubTaskMeta{
	Name:             "ExtractCard",
	EntryPoint:       ExtractCard,
	EnabledByDefault: true,
	Description:      "Extract raw data into tool layer table trello_cards",
}

type TrelloApiCard struct {
	ID                    string        `json:"id"`
	Badges                interface{}   `json:"badges"`
	CheckItemStates       interface{}   `json:"checkItemStates"`
	Closed                bool          `json:"closed"`
	DueComplete           bool          `json:"dueComplete"`
	DateLastActivity      time.Time     `json:"dateLastActivity"`
	Desc                  string        `json:"desc"`
	DescData              interface{}   `json:"descData"`
	Due                   interface{}   `json:"due"`
	DueReminder           interface{}   `json:"dueReminder"`
	Email                 interface{}   `json:"email"`
	IDBoard               string        `json:"idBoard"`
	IDChecklists          []string      `json:"idChecklists"`
	IDList                string        `json:"idList"`
	IDMembers             []string      `json:"idMembers"`
	IDMembersVoted        []string      `json:"idMembersVoted"`
	IDShort               int           `json:"idShort"`
	IDAttachmentCover     string        `json:"idAttachmentCover"`
	Labels                []interface{} `json:"labels"`
	IDLabels              []string      `json:"idLabels"`
	ManualCoverAttachment bool          `json:"manualCoverAttachment"`
	Name                  string        `json:"name"`
	Pos                   float64       `json:"pos"`
	ShortLink             string        `json:"shortLink"`
	ShortUrl              string        `json:"shortUrl"`
	Start                 interface{}   `json:"start"`
	Subscribed            bool          `json:"subscribed"`
	Url                   string        `json:"url"`
	Cover                 interface{}   `json:"cover"`
	IsTemplate            bool          `json:"isTemplate"`
	CardRole              interface{}   `json:"cardRole"`
}

func ExtractCard(taskCtx plugin.SubTaskContext) errors.Error {
	taskData := taskCtx.GetData().(*TrelloTaskData)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TrelloApiParams{
				ConnectionId: taskData.Options.ConnectionId,
				BoardId:      taskData.Options.BoardId,
			},
			Table: RAW_CARD_TABLE,
		},
		Extract: func(resData *api.RawData) ([]interface{}, errors.Error) {
			apiCard := &TrelloApiCard{}
			err := errors.Convert(json.Unmarshal(resData.Data, apiCard))
			if err != nil {
				return nil, err
			}
			return []interface{}{
				&models.TrelloCard{
					ID:               apiCard.ID,
					Name:             apiCard.Name,
					Closed:           apiCard.Closed,
					DueComplete:      apiCard.DueComplete,
					DateLastActivity: apiCard.DateLastActivity,
					IDBoard:          apiCard.IDBoard,
					IDList:           apiCard.IDList,
					IDShort:          apiCard.IDShort,
					Pos:              apiCard.Pos,
					ShortLink:        apiCard.ShortLink,
					ShortUrl:         apiCard.ShortUrl,
					Subscribed:       apiCard.Subscribed,
					Url:              apiCard.Url,
				},
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}
