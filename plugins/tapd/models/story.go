package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdStory struct {
	SourceId        uint64 `gorm:"primaryKey"`
	ID              uint64 `gorm:"primaryKey;type:BIGINT(100)" json:"id"`
	WorkitemTypeID  uint64 `json:"workitem_type_id"`
	EpicKey         string
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	WorkspaceId     uint64     `json:"workspace_id"`
	Creator         string     `json:"creator"`
	Created         *time.Time `json:"created"`
	Modified        *time.Time `json:"modified" gorm:"index"`
	Status          string     `json:"status"`
	Owner           string     `json:"owner"`
	Cc              string     `json:"cc"`
	Begin           *time.Time `json:"begin"`
	Due             *time.Time `json:"due"`
	Size            int        `json:"size"`
	Priority        string     `json:"priority"`
	Developer       string     `json:"developer"`
	IterationID     uint64     `json:"iteration_id"`
	TestFocus       string     `json:"test_focus"`
	Type            string     `json:"type"`
	Source          string     `json:"source"`
	Module          string     `json:"module"`
	Version         string     `json:"version"`
	Completed       *time.Time `json:"completed"`
	CategoryID      uint64     `json:"category_id"`
	Path            string     `json:"path"`
	ParentID        uint64     `json:"parent_id"`
	ChildrenID      string     `json:"children_id"`
	AncestorID      uint64     `json:"ancestor_id"`
	BusinessValue   string     `json:"business_value"`
	Effort          uint64     `json:"effort"`
	EffortCompleted uint64     `json:"effort_completed"`
	Exceed          uint64     `json:"exceed"`
	Remain          uint64     `json:"remain"`
	ReleaseID       uint64     `json:"release_id"`
	Confidential    string     `json:"confidential"`
	TemplatedID     uint64     `json:"templated_id"`
	CreatedFrom     string     `json:"created_from"`
	Feature         string     `json:"feature"`
	Url             string
	common.NoPKModel
}

type TapdStoryApiRes struct {
	ID              string `gorm:"primaryKey" json:"id"`
	WorkitemTypeID  string `json:"workitem_type_id"`
	EpicKey         string
	Name            string `json:"name"`
	Description     string `json:"description"`
	WorkspaceId     string `json:"workspace_id"`
	Creator         string `json:"creator"`
	Created         string `json:"created"`
	Modified        string `json:"modified" gorm:"index"`
	Status          string `json:"status"`
	Owner           string `json:"owner"`
	Cc              string `json:"cc"`
	Begin           string `json:"begin"`
	Due             string `json:"due"`
	Size            string `json:"size"`
	Priority        string `json:"priority"`
	Developer       string `json:"developer"`
	IterationID     string `json:"iteration_id"`
	TestFocus       string `json:"test_focus"`
	Type            string `json:"type"`
	Source          string `json:"source"`
	Module          string `json:"module"`
	Version         string `json:"version"`
	Completed       string `json:"completed"`
	CategoryID      string `json:"category_id"`
	Path            string `json:"path"`
	ParentID        string `json:"parent_id"`
	ChildrenID      string `json:"children_id"`
	AncestorID      string `json:"ancestor_id"`
	BusinessValue   string `json:"business_value"`
	Effort          string `json:"effort"`
	EffortCompleted string `json:"effort_completed"`
	Exceed          string `json:"exceed"`
	Remain          string `json:"remain"`
	ReleaseID       string `json:"release_id"`
	Confidential    string `json:"confidential"`
	TemplatedID     string `json:"templated_id"`
	CreatedFrom     string `json:"created_from"`
	Feature         string `json:"feature"`
}
