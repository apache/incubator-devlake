package archived

import (
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type ZentaoExecution struct {
	ConnectionId   uint64              `gorm:"primaryKey"`
	Id             uint64              `json:"id" gorm:"primaryKey"`
	Project        uint64              `json:"project"`
	Model          string              `json:"model"`
	Type           string              `json:"type"`
	Lifetime       string              `json:"lifetime"`
	Budget         string              `json:"budget"`
	BudgetUnit     string              `json:"budgetUnit"`
	Attribute      string              `json:"attribute"`
	Percent        int                 `json:"percent"`
	Milestone      string              `json:"milestone"`
	Output         string              `json:"output"`
	Auth           string              `json:"auth"`
	Parent         int                 `json:"parent"`
	Path           string              `json:"path"`
	Grade          int                 `json:"grade"`
	Name           string              `json:"name"`
	Code           string              `json:"code"`
	Begin          *helper.Iso8601Time `json:"begin"`
	End            *helper.Iso8601Time `json:"end"`
	RealBegan      *helper.Iso8601Time `json:"realBegan"`
	RealEnd        *helper.Iso8601Time `json:"realEnd"`
	Days           int                 `json:"days"`
	Status         string              `json:"status"`
	SubStatus      string              `json:"subStatus"`
	Pri            string              `json:"pri"`
	Desc           string              `json:"desc"`
	Version        int                 `json:"version"`
	ParentVersion  int                 `json:"parentVersion"`
	PlanDuration   int                 `json:"planDuration"`
	RealDuration   int                 `json:"realDuration"`
	OpenedBy       `json:"openedBy"`
	OpenedDate     *helper.Iso8601Time `json:"openedDate"`
	OpenedVersion  string              `json:"openedVersion"`
	LastEditedBy   `json:"lastEditedBy"`
	LastEditedDate *helper.Iso8601Time `json:"lastEditedDate"`
	ClosedBy       `json:"closedBy"`
	ClosedDate     *helper.Iso8601Time `json:"closedDate"`
	CanceledBy     `json:"canceledBy"`
	CanceledDate   *helper.Iso8601Time `json:"canceledDate"`
	SuspendedDate  *helper.Iso8601Time `json:"suspendedDate"`
	PO             `json:"PO"`
	PM             `json:"PM"`
	QD             `json:"QD"`
	RD             `json:"RD"`
	Team           string `json:"team"`
	Acl            string `json:"acl"`
	//Whitelist      []Whitelist  `json:"whitelist" gorm:"-:all"`
	Order         int          `json:"order"`
	Vision        string       `json:"vision"`
	DisplayCards  int          `json:"displayCards"`
	FluidBoard    string       `json:"fluidBoard"`
	Deleted       bool         `json:"deleted"`
	TotalHours    int          `json:"totalHours"`
	TotalEstimate int          `json:"totalEstimate"`
	TotalConsumed int          `json:"totalConsumed"`
	TotalLeft     int          `json:"totalLeft"`
	ProjectInfo   bool         `json:"projectInfo"`
	Progress      int          `json:"progress"`
	TeamMembers   []TeamMember `json:"teamMembers" gorm:"-:all"`
	Products      []Product    `json:"products" gorm:"-:all"`
	CaseReview    bool         `json:"caseReview"`
	archived.NoPKModel
}

type OpenedBy struct {
	OpenedByID       int    `json:"id"`
	OpenedByAccount  string `json:"account"`
	OpenedByAvatar   string `json:"avatar"`
	OpenedByRealname string `json:"realname"`
}

type LastEditedBy struct {
	LastEditedByID       int    `json:"id"`
	LastEditedByAccount  string `json:"account"`
	LastEditedByAvatar   string `json:"avatar"`
	LastEditedByRealname string `json:"realname"`
}

type ClosedBy struct {
	ClosedByID       int    `json:"id"`
	ClosedByAccount  string `json:"account"`
	ClosedByAvatar   string `json:"avatar"`
	ClosedByRealname string `json:"realname"`
}

type CanceledBy struct {
	CanceledByID       int    `json:"id"`
	CanceledByAccount  string `json:"account"`
	CanceledByAvatar   string `json:"avatar"`
	CanceledByRealname string `json:"realname"`
}

type PO struct {
	PoID       int    `json:"id"`
	PoAccount  string `json:"account"`
	PoAvatar   string `json:"avatar"`
	PoRealname string `json:"realname"`
}

type QD struct {
	ID       int    `json:"id"`
	Account  string `json:"account"`
	Avatar   string `json:"avatar"`
	Realname string `json:"realname"`
}

type RD struct {
	ID       int    `json:"id"`
	Account  string `json:"account"`
	Avatar   string `json:"avatar"`
	Realname string `json:"realname"`
}

type Product struct {
	ID    int           `json:"id"`
	Name  string        `json:"name"`
	Plans []interface{} `json:"plans"`
}

type TeamMember struct {
	ID         int    `json:"id"`
	Root       int    `json:"root"`
	Type       string `json:"type"`
	Account    string `json:"account"`
	Role       string `json:"role"`
	Position   string `json:"position"`
	Limited    string `json:"limited"`
	Join       string `json:"join"`
	Days       int    `json:"days"`
	Hours      int    `json:"hours"`
	Estimate   string `json:"estimate"`
	Consumed   string `json:"consumed"`
	Left       string `json:"left"`
	Order      int    `json:"order"`
	TotalHours int    `json:"totalHours"`
	UserID     int    `json:"userID"`
	Realname   string `json:"realname"`
}

func (ZentaoExecution) TableName() string {
	return "_tool_zentao_executions"
}
