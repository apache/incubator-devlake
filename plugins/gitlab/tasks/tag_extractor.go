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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type GitlabApiTag struct {
	Name      string
	Message   string
	Target    string
	Protected bool
	Release   struct {
		TagName     string
		Description string
	}
}

var ExtractTagMeta = core.SubTaskMeta{
	Name:             "extractApiTag",
	EntryPoint:       ExtractApiTag,
	EnabledByDefault: false,
	Description:      "Extract raw tag data into tool layer table GitlabTag",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE},
}

func ExtractApiTag(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TAG_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			// need to extract 1 kinds of entities here
			results := make([]interface{}, 0, 1)

			// create gitlab commit
			gitlabApiTag := &GitlabApiTag{}
			err := json.Unmarshal(row.Data, gitlabApiTag)
			if err != nil {
				return nil, err
			}
			gitlabTag, err := convertTag(gitlabApiTag)
			if err != nil {
				return nil, err
			}
			gitlabTag.ConnectionId = data.Options.ConnectionId
			results = append(results, gitlabTag)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

// Convert the API response to our DB model instance
func convertTag(tag *GitlabApiTag) (*models.GitlabTag, error) {
	gitlabTag := &models.GitlabTag{
		Name:               tag.Name,
		Message:            tag.Message,
		Target:             tag.Target,
		Protected:          tag.Protected,
		ReleaseDescription: tag.Release.Description,
	}
	return gitlabTag, nil
}
