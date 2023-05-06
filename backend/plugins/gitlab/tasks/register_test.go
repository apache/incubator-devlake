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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/subTaskMetaSorter"
	"testing"
)

func TestGetSubtaskList(t *testing.T) {
	sorter := subTaskMetaSorter.GetDependencySorter(SubTaskMetaList)
	list, err := sorter.Sort()
	if err != nil {
		t.Error(err)
		return
	}
	rawList := []plugin.SubTaskMeta{
		CollectApiIssuesMeta,
		ExtractApiIssuesMeta,
		CollectApiMergeRequestsMeta,
		ExtractApiMergeRequestsMeta,
		CollectApiMergeRequestDetailsMeta,
		CollectApiMrNotesMeta,
		ExtractApiMrNotesMeta,
		CollectApiMrCommitsMeta,
		ExtractApiMrCommitsMeta,
		CollectApiPipelinesMeta,
		ExtractApiPipelinesMeta,
		CollectApiPipelineDetailsMeta,
		ExtractApiPipelineDetailsMeta,
		CollectApiJobsMeta,
		ExtractApiJobsMeta,
		EnrichMergeRequestsMeta,
		CollectAccountsMeta,
		ExtractAccountsMeta,
		ConvertAccountsMeta,
		ConvertProjectMeta,
		ConvertApiMergeRequestsMeta,
		ConvertMrCommentMeta,
		ConvertApiMrCommitsMeta,
		ConvertIssuesMeta,
		ConvertIssueLabelsMeta,
		ConvertMrLabelsMeta,
		ConvertCommitsMeta,
		ConvertPipelineMeta,
		ConvertPipelineCommitMeta,
		ConvertJobMeta,
		CollectApiCommitsMeta,
		ExtractApiCommitsMeta,
		ExtractApiMergeRequestDetailsMeta,
		CollectTagMeta,
		ExtractTagMeta,
	}
	if len(list) != len(rawList) {
		t.Errorf("get wrong length with raw list, raw = %d list = %d", len(rawList), len(list))
		return
	}
	for index, item := range rawList {
		if item.Name != list[index].Name {
			t.Errorf("get wrong item name in given index, index[%d] rawitem[%s] item[]%s]",
				index, item.Name, list[index].Name)
			return
		}
	}
}
