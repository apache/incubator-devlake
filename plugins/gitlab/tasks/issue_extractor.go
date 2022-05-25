package tasks

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/plugins/core"

	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractApiIssuesMeta = core.SubTaskMeta{
	Name:             "extractApiIssues",
	EntryPoint:       ExtractApiIssues,
	EnabledByDefault: true,
	Description:      "Extract raw Issues data into tool layer table gitlab_issues",
}

type IssuesResponse struct {
	ProjectId    int `json:"id"`
	Milestone    struct {
		Due_date 	string
		Project_id 	int
		State 		string
		Description string
		Iid 		int
		Id 			int
		Title 		string
		CreatedAt *core.Iso8601Time
		UpdatedAt *core.Iso8601Time
	}
	Author struct{
		State 		string
		WebUrl 		string
		AvatarUrl 	string
		Username 	string
		Id 			int
		Name 		string
	}
	Description 	string
	State      		string
	Iid 			int
	Assignees []struct {
		AvatarUrl 	string
		WebUrl 		string
		State 		string
		Username 	string
		Id 			int
		Name 		string
	}
	Assignee *struct {
		AvatarUrl 	string
		WebUrl 		string
		State 		string
		Username 	string
		Id 			int
		Name 		string
	}
	Type 				string
	Labels 				[]string `json:"labels"`
	UpVotes 			int
	DownVotes 			int
	MergeRequestsCount 	int
	Id 					int
	Title       		string
	GitlabUpdatedAt core.Iso8601Time  `json:"updated_at"`
	GitlabCreatedAt core.Iso8601Time  `json:"created_at"`
	GitlabClosedAt  *core.Iso8601Time `json:"closed_at"`
	ClosedBy struct{
		State 		string
		WebUrl 		string
		AvatarUrl	string
		Username 	string
		Id 			int
		Name 		string
	}
	UserNotesCount int
	DueDate *core.Iso8601Time
	WebUrl string	`json:"web_url"`
	References struct {
		Short 		string
		Relative 	string
		Full 		string
	}
	TimeStats struct {
		TimeEstimate 		int
		TotalTimeSpent 		int
		HumanTimeEstimate 	string
		HumanTotalTimeSpent string
	}
	HasTasks 		bool
	TaskStatus 		string
	Confidential 	bool
	DiscussionLocked bool
	IssueType 		string
	Serverity 		string
	Links struct {
		Self 		string 	`json:"url"`
		Notes 		string
		AwardEmoji 	string
		Project 	string
	}
	TaskCompletionStatus struct {
		Count 			int
		CompletedCount 	int
	}

}

func ExtractApiIssues(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GitlabTaskData)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: GitlabApiParams{
				ProjectId: data.Options.ProjectId,
			},
			/*
				Table store raw data
			*/
			Table: RAW_ISSUE_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			body := &IssuesResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}
			// need to extract 2 kinds of entities here
			if body.ProjectId == 0 {
				return nil, nil
			}
			////If this is a pr, ignore
			//if body.PullRequest.Url != "" {
			//	return nil, nil
			//}
			//If this is not Issue, ignore
			if body.IssueType != "ISSUE" {
				return nil, nil
			}
			results := make([]interface{}, 0, 2)
			gitlabIssue, err := convertGitlabIssue(body, data.Options.ProjectId)
			if err != nil {
				return nil, err
			}

			for _, label := range body.Labels {
				results = append(results, &models.GitlabIssueLabel{
					IssueId:   gitlabIssue.GitlabId,
					LabelName: label,
				})

			}
			results = append(results, gitlabIssue)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
func convertGitlabIssue(issue *IssuesResponse, projectId int) (*models.GitlabIssue, error) {
	gitlabIssue := &models.GitlabIssue{
		GitlabId:        issue.Id,
		ProjectId:       projectId,
		Number:          issue.Iid,
		State:           issue.State,
		Title:           issue.Title,
		Body:            issue.Description,
		Url:             issue.Links.Self,
		ClosedAt:        core.Iso8601TimeToTime(issue.GitlabClosedAt),
		GitlabCreatedAt: issue.GitlabCreatedAt.ToTime(),
		GitlabUpdatedAt: issue.GitlabUpdatedAt.ToTime(),
	}

	if issue.Assignee != nil {
		gitlabIssue.AssigneeId = issue.Assignee.Id
		gitlabIssue.AssigneeName = issue.Assignee.Username
	}
	if issue.GitlabClosedAt != nil {
		gitlabIssue.LeadTimeMinutes = uint(issue.GitlabClosedAt.ToTime().Sub(issue.GitlabCreatedAt.ToTime()).Minutes())
	}

	return gitlabIssue, nil
}
