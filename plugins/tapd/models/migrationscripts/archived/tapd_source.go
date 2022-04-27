package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

type TapdSource struct {
	common.Model
	Name                       string      `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	Endpoint                   string      `gorm:"type:varchar(255)"`
	BasicAuthEncoded           string      `gorm:"type:varchar(255)"`
	EpicKeyField               string      `gorm:"type:varchar(50);" json:"epicKeyField"`
	StoryPointField            string      `gorm:"type:varchar(50);" json:"storyPointField"`
	RemotelinkCommitShaPattern string      `gorm:"type:varchar(255);comment='golang regexp, the first group will be recognized as commit sha, ref https://github.com/google/re2/wiki/Syntax'" json:"remotelinkCommitShaPattern"`
	Proxy                      string      `gorm:"type:varchar(255)"`
	RateLimit                  models.Ints `comment:"api request rate limt per second"`
	//CompanyId                  models.Uint64s `json:"companyId" validate:"required"`
	WorkspaceID models.Uint64s `json:"workspaceId" validate:"required"`
}

type TapdSourceDetail struct {
	TapdSource
}

func (TapdSource) TableName() string {
	return "_tool_tapd_sources"
}
