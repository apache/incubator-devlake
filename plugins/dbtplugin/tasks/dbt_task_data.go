package tasks

// import (
// 	"github.com/merico-dev/lake/plugins/helper"
// )

type DbtOptions struct {
	SelectedModels string   `json:"selectedModels"`
	Tasks          []string `json:"tasks,omitempty"`
}

type DbtTaskData struct {
	Options *DbtOptions
}
