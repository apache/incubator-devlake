package archived

import "github.com/apache/incubator-devlake/plugins/helper"

type TapdWorkitemType struct {
	ConnectionId   uint64          `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Id             uint64          `gorm:"primaryKey;type:BIGINT  NOT NULL" json:"id,string"`
	WorkspaceId    uint64          `json:"workspace_id,string"`
	EntityType     string          `gorm:"type:varchar(255)" json:"entity_type"`
	Name           string          `gorm:"type:varchar(255)" json:"name"`
	EnglishName    string          `gorm:"type:varchar(255)" json:"english_name"`
	Status         string          `gorm:"type:varchar(255)" json:"status"`
	Color          string          `gorm:"type:varchar(255)" json:"color"`
	WorkflowID     uint64          `json:"workflow_id"`
	Icon           string          `json:"icon"`
	IconSmall      string          `json:"icon_small"`
	Creator        string          `gorm:"type:varchar(255)" json:"creator"`
	Created        *helper.CSTTime `json:"created"`
	ModifiedBy     string          `gorm:"type:varchar(255)" json:"modified_by"`
	Modified       *helper.CSTTime `json:"modified"`
	IconViper      string          `json:"icon_viper"`
	IconSmallViper string          `json:"icon_small_viper"`
}

func (TapdWorkitemType) TableName() string {
	return "_tool_tapd_workitem_type"
}
