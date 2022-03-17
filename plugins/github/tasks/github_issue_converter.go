package tasks

import (
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"reflect"
	"strconv"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
)

func ConvertIssues(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	issue := &githubModels.GithubIssue{}
	cursor, err := db.Model(issue).
		Where("_raw_data_params = ?", data.Options.ParamString).
		Rows()

	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubIssue{})
	userIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubUser{})
	boardIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubRepo{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		Ctx:          taskCtx,
		InputRowType: reflect.TypeOf(githubModels.GithubIssue{}),
		Input:        cursor,
		BatchSelectors: map[reflect.Type]helper.BatchSelector{
			reflect.TypeOf(&ticket.Issue{}): {
				Query: "_raw_data_params = ?",
				Parameters: []interface{}{
					data.Options.ParamString,
				},
			},
			reflect.TypeOf(&ticket.BoardIssue{}): {
				Query: "_raw_data_params = ?",
				Parameters: []interface{}{
					data.Options.ParamString,
				},
			},
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			issue := inputRow.(*githubModels.GithubIssue)
			domainIssue := &ticket.Issue{
				DomainEntity:    domainlayer.DomainEntity{Id: issueIdGen.Generate(issue.GithubId)},
				Key:             strconv.Itoa(issue.Number),
				Title:           issue.Title,
				Summary:         issue.Body,
				Priority:        issue.Priority,
				Type:            issue.Type,
				AssigneeId:      userIdGen.Generate(issue.AssigneeId),
				AssigneeName:    issue.AssigneeName,
				LeadTimeMinutes: issue.LeadTimeMinutes,
				CreatedDate:     &issue.GithubCreatedAt,
				UpdatedDate:     &issue.GithubUpdatedAt,
				ResolutionDate:  issue.ClosedAt,
				Severity:        issue.Severity,
				Component:       issue.Component,
			}
			if issue.State == "closed" {
				domainIssue.Status = ticket.DONE
			} else {
				domainIssue.Status = ticket.TODO
			}
			domainIssue.RawDataOrigin = issue.RawDataOrigin

			boardIssue := &ticket.BoardIssue{
				BoardId: boardIdGen.Generate(repoId),
				IssueId: domainIssue.Id,
			}
			boardIssue.RawDataOrigin = issue.RawDataOrigin

			return []interface{}{
				domainIssue,
				boardIssue,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
