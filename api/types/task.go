package types

type CreateTask struct {
	// Plugin name
	Plugin string `json:"plugin" binding:"required" default:"Jira" validate:"required"`
	// Options for the plugin task to be triggered
	Options map[string]interface{} `json:"options" binding:"required"`
}
