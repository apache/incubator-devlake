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
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/plugins/core"
)

type RefdiffOptions struct {
	RepoId string
	Tasks  []string `json:"tasks,omitempty"`
	Pairs  []RefPair

	TagsPattern string // The Pattern to match from all tags
	TagsLimit   int    // How many tags be matched should be used.
	TagsOrder   string // The Rule to Order the tag list
}

type RefdiffTaskData struct {
	Options *RefdiffOptions
	Since   *time.Time
}

type RefPair struct {
	NewRef string
	OldRef string
}

type Refs []code.Ref
type RefsAlphabetically Refs
type RefsReverseAlphabetically Refs

func (rs Refs) Len() int {
	return len(rs)
}

func (rs RefsAlphabetically) Len() int {
	return len(rs)
}

func (rs RefsAlphabetically) Less(i, j int) bool {
	return rs[i].Id < rs[j].Id
}

func (rs RefsAlphabetically) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs RefsReverseAlphabetically) Len() int {
	return len(rs)
}

func (rs RefsReverseAlphabetically) Less(i, j int) bool {
	return rs[i].Id > rs[j].Id
}

func (rs RefsReverseAlphabetically) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

// Calculate the TagPattern order by tagsOrder and return the Refs
func CaculateTagPattern(taskCtx core.SubTaskContext) (Refs, error) {
	rs := Refs{}
	data := taskCtx.GetData().(*RefdiffTaskData)
	tagsPattern := data.Options.TagsPattern
	tagsOrder := data.Options.TagsOrder
	db := taskCtx.GetDb()

	// caculate Pattern part
	if data.Options.TagsPattern == "" || data.Options.TagsLimit <= 1 {
		return rs, nil
	}
	rows, err := db.Model(&code.Ref{}).Order("created_date desc").Rows()
	if err != nil {
		return rs, err
	}
	defer rows.Next()
	r, err := regexp.Compile(tagsPattern)
	if err != nil {
		return rs, fmt.Errorf("unable to parse: %s\r\n%s", tagsPattern, err.Error())
	}
	for rows.Next() {
		var ref code.Ref
		err = db.ScanRows(rows, &ref)
		if err != nil {
			return rs, err
		}

		if ok := r.Match([]byte(ref.Id)); ok {
			rs = append(rs, ref)
		}
	}
	switch tagsOrder {
	case "alphabetically":
		sort.Sort(RefsAlphabetically(rs))
	case "reverse alphabetically":
		sort.Sort(RefsReverseAlphabetically(rs))
	default:
		break
	}
	return rs, nil
}
