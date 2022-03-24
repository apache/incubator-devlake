package apiv2models

import "time"

type Input struct {
	IssueId    uint64    `json:"issue_id"`
	UpdateTime time.Time `json:"update_time"`
}
