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

package models

import (
	"encoding/json"
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
)

type ZentaoStoryRes struct {
	AllFeilds        map[string]interface{} `json:"-"`
	ID               int64                  `json:"id"`
	Vision           string                 `json:"vision"`
	Parent           int64                  `json:"parent"`
	Product          int64                  `json:"product"`
	Branch           int                    `json:"branch"`
	Module           int                    `json:"module"`
	Plan             string                 `json:"plan"`
	Source           string                 `json:"source"`
	SourceNote       string                 `json:"sourceNote"`
	FromBug          int                    `json:"fromBug"`
	Feedback         int                    `json:"feedback"`
	Title            string                 `json:"title"`
	Keywords         string                 `json:"keywords"`
	Type             string                 `json:"type"`
	Category         string                 `json:"category"`
	Pri              int                    `json:"pri"`
	Estimate         float64                `json:"estimate"`
	Status           string                 `json:"status"`
	SubStatus        string                 `json:"subStatus"`
	Color            string                 `json:"color"`
	Stage            string                 `json:"stage"`
	Mailto           []interface{}          `json:"mailto"`
	Lib              int                    `json:"lib"`
	FromStory        int64                  `json:"fromStory"`
	FromVersion      int                    `json:"fromVersion"`
	OpenedBy         *ApiAccount            `json:"openedBy"`
	OpenedDate       *common.Iso8601Time    `json:"openedDate"`
	AssignedTo       *ApiAccount            `json:"assignedTo"`
	AssignedDate     *common.Iso8601Time    `json:"assignedDate"`
	ApprovedDate     *common.Iso8601Time    `json:"approvedDate"`
	LastEditedBy     *ApiAccount            `json:"lastEditedBy"`
	LastEditedDate   *common.Iso8601Time    `json:"lastEditedDate"`
	ChangedBy        string                 `json:"changedBy"`
	ChangedDate      *common.Iso8601Time    `json:"changedDate"`
	ReviewedBy       *ApiAccount            `json:"reviewedBy"`
	ReviewedDate     *common.Iso8601Time    `json:"reviewedDate"`
	ClosedBy         *ApiAccount            `json:"closedBy"`
	ClosedDate       *common.Iso8601Time    `json:"closedDate"`
	ClosedReason     string                 `json:"closedReason"`
	ActivatedDate    *common.Iso8601Time    `json:"activatedDate"`
	ToBug            int                    `json:"toBug"`
	ChildStories     string                 `json:"childStories"`
	LinkStories      string                 `json:"linkStories"`
	LinkRequirements string                 `json:"linkRequirements"`
	DuplicateStory   int64                  `json:"duplicateStory"`
	Version          int                    `json:"version"`
	StoryChanged     string                 `json:"storyChanged"`
	FeedbackBy       string                 `json:"feedbackBy"`
	NotifyEmail      string                 `json:"notifyEmail"`
	URChanged        string                 `json:"URChanged"`
	Deleted          bool                   `json:"deleted"`
	PriOrder         *common.StringFloat64  `json:"priOrder"`
	PlanTitle        string                 `json:"planTitle"`
	ProductStatus    string                 `json:"productStatus"`
}

func (i *ZentaoStoryRes) SetAllFeilds(raw json.RawMessage) error {
	var allFeilds map[string]interface{}
	if err := json.Unmarshal(raw, &allFeilds); err != nil {
		return err
	}
	i.AllFeilds = allFeilds
	return nil
}

type ZentaoStory struct {
	common.NoPKModel
	ConnectionId uint64  `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ID           int64   `json:"id" gorm:"primaryKey;type:BIGINT  NOT NULL;autoIncrement:false" `
	Product      int64   `json:"product"`
	Branch       int     `json:"branch"`
	Version      int     `json:"version"`
	OrderIn      int     `json:"order"`
	Vision       string  `json:"vision"`
	Parent       int64   `json:"parent"`
	Module       int     `json:"module"`
	Plan         string  `json:"plan"`
	Source       string  `json:"source"`
	SourceNote   string  `json:"sourceNote"`
	FromBug      int     `json:"fromBug"`
	Feedback     int     `json:"feedback"`
	Title        string  `json:"title"`
	Keywords     string  `json:"keywords"`
	Type         string  `json:"type"`
	Category     string  `json:"category"`
	Pri          int     `json:"pri"`
	Estimate     float64 `json:"estimate"`
	Status       string  `json:"status"`
	SubStatus    string  `json:"subStatus"`
	Color        string  `json:"color"`
	Stage        string  `json:"stage"`
	//Mailto           []interface{} `json:"mailto"`
	Lib              int   `json:"lib"`
	FromStory        int64 `json:"fromStory"`
	FromVersion      int   `json:"fromVersion"`
	OpenedById       int64
	OpenedByName     string
	OpenedDate       *common.Iso8601Time `json:"openedDate"`
	AssignedToId     int64
	AssignedToName   string
	AssignedDate     *common.Iso8601Time `json:"assignedDate"`
	ApprovedDate     *common.Iso8601Time `json:"approvedDate"`
	LastEditedId     int64
	LastEditedDate   *common.Iso8601Time `json:"lastEditedDate"`
	ChangedDate      *common.Iso8601Time `json:"changedDate"`
	ReviewedById     int64               `json:"reviewedBy"`
	ReviewedDate     *common.Iso8601Time `json:"reviewedDate"`
	ClosedId         int64
	ClosedDate       *common.Iso8601Time `json:"closedDate"`
	ClosedReason     string              `json:"closedReason"`
	ActivatedDate    *common.Iso8601Time `json:"activatedDate"`
	ToBug            int                 `json:"toBug"`
	ChildStories     string              `json:"childStories"`
	LinkStories      string              `json:"linkStories"`
	LinkRequirements string              `json:"linkRequirements"`
	DuplicateStory   int64               `json:"duplicateStory"`
	StoryChanged     string              `json:"storyChanged"`
	FeedbackBy       string              `json:"feedbackBy"`
	NotifyEmail      string              `json:"notifyEmail"`
	URChanged        string              `json:"URChanged"`
	Deleted          bool                `json:"deleted"`
	PriOrder         string              `json:"priOrder"`
	PlanTitle        string              `json:"planTitle"`
	Url              string              `json:"url"`
	StdStatus        string              `json:"stdStatus" gorm:"type:varchar(20)"`
	StdType          string              `json:"stdType" gorm:"type:varchar(20)"`
	DueDate          *time.Time          `json:"dueDate"`
}

func (ZentaoStory) TableName() string {
	return "_tool_zentao_stories"
}
