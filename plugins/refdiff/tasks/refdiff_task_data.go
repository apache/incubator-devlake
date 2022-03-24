package tasks

import (
	"time"
)

type RefdiffOptions struct {
	RepoId string
	Tasks  []string `json:"tasks,omitempty"`
	Pairs  []RefPair
}

type RefdiffTaskData struct {
	Options *RefdiffOptions
	Since   *time.Time
}

type RefPair struct {
	NewRef string
	OldRef string
}
