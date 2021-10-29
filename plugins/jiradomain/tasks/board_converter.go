package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/domainlayer/models/base"
	"github.com/merico-dev/lake/plugins/domainlayer/models/ticket"
	"github.com/merico-dev/lake/plugins/domainlayer/okgen"
	jiraModels "github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm/clause"
)

func ConvertBoard(sourceId uint64, boardId uint64) error {
	jiraBoard := &jiraModels.JiraBoard{}

	err := lakeModels.Db.First(jiraBoard, "source_id = ? AND board_id = ?", sourceId, boardId).Error
	if err != nil {
		return err
	}

	board := &ticket.Board{
		DomainEntity: base.DomainEntity{
			OriginKey: okgen.NewOriginKeyGenerator(jiraBoard).Generate(jiraBoard.SourceId, boardId),
		},
		Name: jiraBoard.Name,
		Url:  jiraBoard.Self,
	}

	return lakeModels.Db.Clauses(clause.OnConflict{UpdateAll: true}).Create(board).Error
}
