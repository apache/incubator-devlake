package tasks

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"reflect"
	"time"
)

type BugChangelogItemResult struct {
	SourceId          uint64    `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	WorkspaceId       uint64    `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID                uint64    `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id"`
	BugID             uint64    `json:"bug_id"`
	Author            string    `json:"author"`
	Field             string    `json:"field"`
	OldValue          string    `json:"old_value"`
	NewValue          string    `json:"new_value"`
	Memo              string    `json:"memo"`
	Created           time.Time `json:"created"`
	ChangelogId       uint64    `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL"`
	ValueBeforeParsed string    `json:"value_before"`
	ValueAfterParsed  string    `json:"value_after"`
	IterationIdFrom   uint64
	IterationIdTo     uint64
	common.NoPKModel
}

func ConvertBugChangelog(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("convert changelog :%d", data.Options.WorkspaceId)
	clIdGen := didgen.NewDomainIdGenerator(&models.TapdBugChangelog{})

	cursor, err := db.Table("_tool_tapd_bug_changelog_items").
		Joins("left join _tool_tapd_bug_changelogs tc on tc.id = _tool_tapd_bug_changelog_items.changelog_id ").
		Where("tc.source_id = ? AND tc.workspace_id = ?", data.Source.ID, data.Options.WorkspaceId).
		Select("tc.created, tc.id, tc.workspace_id, tc.bug_id, tc.author, _tool_tapd_bug_changelog_items.*").
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	changelogToHistoryConverter := NewChangelogToHistoryConverter(taskCtx)
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId:   data.Source.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_BUG_CHANGELOG_TABLE,
		},
		InputRowType: reflect.TypeOf(BugChangelogItemResult{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			cl := inputRow.(*BugChangelogItemResult)
			domainCl := &ticket.Changelog{
				DomainEntity: domainlayer.DomainEntity{
					Id: clIdGen.Generate(data.Source.ID, cl.ID, cl.Field),
				},
				IssueId:     IssueIdGen.Generate(data.Source.ID, cl.BugID),
				AuthorId:    UserIdGen.Generate(data.Source.ID, data.Options.WorkspaceId, cl.Author),
				AuthorName:  cl.Author,
				FieldId:     cl.Field,
				FieldName:   cl.Field,
				From:        cl.ValueBeforeParsed,
				To:          cl.ValueAfterParsed,
				CreatedDate: cl.Created,
			}

			changelogTmp := &models.ChangelogTmp{
				Id:              cl.ID,
				IssueId:         cl.BugID,
				AuthorId:        cl.Author,
				AuthorName:      cl.Author,
				FieldId:         cl.Field,
				FieldName:       cl.Field,
				From:            cl.ValueBeforeParsed,
				To:              cl.ValueBeforeParsed,
				IterationIdFrom: cl.IterationIdFrom,
				IterationIdTo:   cl.IterationIdTo,
				CreatedDate:     cl.Created,
			}

			changelogToHistoryConverter.FeedIn(data.Source.ID, *changelogTmp)

			return []interface{}{
				domainCl,
			}, nil
		},
	})
	if err != nil {
		logger.Info(err.Error())
		return err
	}

	err = converter.Execute()
	if err != nil {
		return err
	}
	return changelogToHistoryConverter.UpdateSprintIssue()
}

var ConvertBugChangelogMeta = core.SubTaskMeta{
	Name:             "convertBugChangelog",
	EntryPoint:       ConvertBugChangelog,
	EnabledByDefault: true,
	Description:      "convert Tapd bug changelog",
}
