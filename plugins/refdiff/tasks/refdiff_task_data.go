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
	"strings"
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"gorm.io/gorm"
)

type RefdiffOptions struct {
	RepoId string
	Tasks  []string `json:"tasks,omitempty"`
	Pairs  []RefPair

	TagsPattern string     // The Pattern to match from all tags
	TagsLimit   int        // How many tags be matched should be used.
	TagsOrder   string     // The Rule to Order the tag list
	TagsRefs    []code.Ref // Caculate out by TagsPattern TagsLimit TagsOrder from db
}

type RefdiffTaskData struct {
	Options *RefdiffOptions
	Since   *time.Time
}

type RefPair struct {
	NewRef string
	OldRef string
}
type RefCommitPair [4]string
type RefPairList [2]string
type RefCommitPairs []RefCommitPair
type RefPairLists []RefPairList

type Refs []code.Ref
type RefsAlphabetically Refs
type RefsReverseAlphabetically Refs
type RefsSemver Refs
type RefsReverseSemver Refs

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

func (rs RefsSemver) Len() int {
	return len(rs)
}

func (rs RefsSemver) Less(i, j int) bool {
	parti := strings.Split(rs[i].Name, ".")
	partj := strings.Split(rs[j].Name, ".")

	for k := 0; k < len(partj); k++ {
		if k >= len(parti) {
			return true
		}

		if len(parti[k]) < len(partj[k]) {
			return true
		}
		if len(parti[k]) > len(partj[k]) {
			return false
		}

		if parti[k] < partj[k] {
			return true
		}
		if parti[k] > partj[k] {
			return false
		}
	}
	return false
}

func (rs RefsSemver) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs RefsReverseSemver) Len() int {
	return len(rs)
}

func (rs RefsReverseSemver) Less(i, j int) bool {
	return RefsSemver(rs).Less(j, i)
}

func (rs RefsReverseSemver) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

// Calculate the TagPattern order by tagsOrder and return the Refs
func CaculateTagPattern(db *gorm.DB, tagsPattern string, tagsLimit int, tagsOrder string) (Refs, error) {
	rs := Refs{}

	// caculate Pattern part
	if tagsPattern == "" || tagsLimit <= 1 {
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

		if ok := r.Match([]byte(ref.Name)); ok {
			rs = append(rs, ref)
		}
	}
	switch tagsOrder {
	case "alphabetically":
		sort.Sort(RefsAlphabetically(rs))
	case "reverse alphabetically":
		sort.Sort(RefsReverseAlphabetically(rs))
	case "semver":
		sort.Sort(RefsSemver(rs))
	case "reverse semver":
		sort.Sort(RefsReverseSemver(rs))
	default:
		break
	}

	if tagsLimit < rs.Len() {
		rs = rs[:tagsLimit]
	}

	return rs, nil
}
