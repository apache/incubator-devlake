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

package api

type blueprintOutput struct {
	ID   int `json:"id" example:"3"`
	Plan [][]struct {
		Plugin   string   `json:"plugin" example:"jira"`
		Subtasks []string `json:"subtasks" example:"collectStatus,extractStatus,collectProjects,extractProjects,collectBoard,extractBoard,collectIssueTypes,extractIssueType,collectIssues,extractIssues,collectIssueChangelogs,extractIssueChangelogs,collectAccounts,collectWorklogs,extractWorklogs,collectRemotelinks,extractRemotelinks,collectSprints,extractSprints,convertBoard,convertIssues,convertWorklogs,convertIssueChangelogs,convertSprints,convertSprintIssues,convertIssueCommits,extractAccounts,convertAccounts,collectEpics,extractEpics"`
		Options  struct {
			BoardID             int `json:"boardId" example:"8"`
			ConnectionID        int `json:"connectionId" example:"1"`
			TransformationRules struct {
				EpicKeyField               string `json:"epicKeyField" example:"customfield_10014"`
				StoryPointField            string `json:"storyPointField" example:"customfield_10024"`
				RemotelinkCommitShaPattern string `json:"remotelinkCommitShaPattern" example:"/commit/([0-9a-f]{40})$"`
				TypeMappings               struct {
					Bug struct {
						StandardType string `json:"standardType" example:"Bug"`
					} `json:"Bug"`
					Incident struct {
						StandardType string `json:"standardType" example:"Incident"`
					} `json:"Incident"`
					Task struct {
						StandardType string `json:"standardType" example:"Requirement"`
					} `json:"Task"`
				} `json:"typeMappings"`
			} `json:"transformationRules"`
		} `json:"options"`
	} `json:"plan"`
	blueprintInput
}
type blueprintInput struct {
	Name     string `json:"name"`
	Settings struct {
		Version     string `json:"version" example:"1.0.0"`
		Connections []struct {
			Plugin       string `json:"plugin" example:"jira"`
			ConnectionID int    `json:"connectionId" example:"1"`
			Scope        []struct {
				Transformation struct {
					EpicKeyField string `json:"epicKeyField" example:"customfield_10014"`
					TypeMappings struct {
						Bug struct {
							StandardType string `json:"standardType" example:"Bug"`
						} `json:"Bug"`
						Incident struct {
							StandardType string `json:"standardType" example:"Incident"`
						} `json:"Incident"`
						Task struct {
							StandardType string `json:"standardType" example:"Requirement"`
						} `json:"Task"`
					} `json:"typeMappings"`
					StoryPointField            string `json:"storyPointField" example:"customfield_10024"`
					RemotelinkCommitShaPattern string `json:"remotelinkCommitShaPattern" example:"/commit/([0-9a-f]{40})$"`
					BugTags                    []struct {
						Self             string `json:"self" example:"https://merico.atlassian.net/rest/api/2/issuetype/10004"`
						ID               string `json:"id" example:"10004"`
						Description      string `json:"description" example:"测试发现的系统缺陷。"`
						IconURL          string `json:"iconUrl" example:"https://merico.atlassian.net/rest/api/2/universal_avatar/view/type/issuetype/avatar/10303?size=medium"`
						Name             string `json:"name" example:"Bug"`
						UntranslatedName string `json:"untranslatedName" example:"Bug"`
						Subtask          bool   `json:"subtask" example:"false"`
						AvatarID         int    `json:"avatarId" example:"10303"`
						HierarchyLevel   int    `json:"hierarchyLevel" example:"0"`
						Key              int    `json:"key" example:"15"`
						Title            string `json:"title" example:"Bug"`
						Value            string `json:"value" example:"Bug"`
						Type             string `json:"type" example:"string"`
					} `json:"bugTags"`
					IncidentTags []struct {
						Self             string `json:"self" example:"https://merico.atlassian.net/rest/api/2/issuetype/10040"`
						ID               string `json:"id" example:"10040"`
						Description      string `json:"description" example:"For system outages or incidents. Created by Jira Service Desk."`
						IconURL          string `json:"iconUrl" example:"https://merico.atlassian.net/rest/api/2/universal_avatar/view/type/issuetype/avatar/10553?size=medium"`
						Name             string `json:"name" example:"Incident"`
						UntranslatedName string `json:"untranslatedName" example:"Incident"`
						Subtask          bool   `json:"subtask" example:"false"`
						AvatarID         int    `json:"avatarId" example:"10553"`
						HierarchyLevel   int    `json:"hierarchyLevel" example:"0"`
						Key              int    `json:"key" example:"6"`
						Title            string `json:"title" example:"Incident"`
						Value            string `json:"value" example:"Incident"`
						Type             string `json:"type" example:"string"`
					} `json:"incidentTags"`
					RequirementTags []struct {
						Self             string `json:"self" example:"https://merico.atlassian.net/rest/api/2/issuetype/10074"`
						ID               string `json:"id" example:"10074"`
						Description      string `json:"description" example:"Tasks track small, distinct pieces of work."`
						IconURL          string `json:"iconUrl" example:"https://merico.atlassian.net/rest/api/2/universal_avatar/view/type/issuetype/avatar/10318?size=medium"`
						Name             string `json:"name" example:"Task"`
						UntranslatedName string `json:"untranslatedName" example:"Task"`
						Subtask          bool   `json:"subtask" example:"false"`
						AvatarID         int    `json:"avatarId" example:"10318"`
						HierarchyLevel   int    `json:"hierarchyLevel" example:"0"`
						Scope            struct {
							Type    string `json:"type" example:"PROJECT"`
							Project struct {
								ID string `json:"id" example:"10033"`
							} `json:"project"`
						} `json:"scope"`
						Key   int    `json:"key" example:"4"`
						Title string `json:"title" example:"Task"`
						Value string `json:"value" example:"Task"`
						Type  string `json:"type" example:"string"`
					} `json:"requirementTags"`
				} `json:"transformation"`
				Options struct {
					BoardID int `json:"boardId" example:"8"`
				} `json:"options"`
				Entities []string `json:"entities" example:"TICKET,CROSS"`
			} `json:"scope"`
		} `json:"connections"`
	} `json:"settings"`
	CronConfig string `json:"cronConfig" example:"0 0 * * *"`
	Enable     bool   `json:"enable" example:"true"`
	Mode       string `json:"mode" example:"NORMAL"`
	IsManual   bool   `json:"isManual" example:"true"`
}
