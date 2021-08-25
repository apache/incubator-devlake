package types

type CreateSource struct {
	Plugin  string                 `json:"plugin" binding:"required" default:"Jira" validate:"required"`
	Name    string                 `json:"name"`
	Options map[string]interface{} `json:"options" binding:"required"`
}
