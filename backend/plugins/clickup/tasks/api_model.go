package tasks

import (
	"encoding/json"
	"strconv"
	"time"
)

type TaskTimeInStatus struct {
	CurrentStatus struct {
		Status string
		Color string
		TotalTime struct {
			ByMinute int `json:"by_minute"`
			Since string `json:"since"`
		} `json:"total_time"`
	} `json:"current_status"`
	StatusHistory []struct {
		Status string
		Color string
		OrderIndex int
		TotalTime struct {
			ByMinute int `json:"by_minute"`
			Since string `json:"since"`
		} `json:"total_time"`
	} `json:"status_history"`
}

type TaskTimeInStatusEnvelope struct {
	TaskId string
	Data TaskTimeInStatus
}

type Task struct {
	Id           string
	CustomId     *string           `json:"custom_id"`
	CustomFields []TaskCustomField `json:"custom_fields"`
	Name         string
	List         TaskList
	Points       float64 `json:"points"`
	Url          string
	TextContent  string `json:"text_content"`
	Priority     TaskPriority
	Parent       *string
	Description  string
	Status       TaskStatus
	Creator      User
	Assignees    []User
	Watchers     []User
	Project      TaskProject
	DateCreated  string `json:"date_created"`
	DateUpdated  string `json:"date_updated"`
	DateDone     string `json:"date_done"`
	DateClosed   string `json:"date_closed"`
	DueDate      string `json:"due_date"`
	StartDate    string `json:"start_date"`
	TimeSpent    string `json:"time_spent"`
	Subtasks     []Subtask
	Space        TaskSpace
}

type TaskPriority struct {
	Priority string
}
type TaskList struct {
	 Id string
}

type TaskCustomField struct {
	Id         string           `json:"id"`
	Name       string           `json:"name"`
	Type       string           `json:"type"`
	Value      interface{}      `json:"value"`
	TypeConfig *json.RawMessage `json:"type_config"`
}

type TaskCustomFieldDropDown struct {
	Default     int                             `json:"default"`
	Placeholder *string                         `json:"placeholder"`
	Options     []TypeCustomFieldDropDownOption `json:"options"`
}

type TypeCustomFieldDropDownOption struct {
	Name string `json:"name"`
}

type TaskSpace struct {
	Id string
}

type Comment struct {
	Id          string
	CommentText string `json:"comment_text"`
	User        User
	Reactions   []Reaction
	Date        string `json:"date"`
}

func (c Comment) DateFormatted(format string) string {
	int, err := strconv.ParseInt(c.Date, 10, 64)
	if err != nil {
		return "nan"
	}
	time := time.Unix(int, 0)
	return time.Format(format)
}

type Comments struct {
	Comments []Comment
}

type Reaction struct {
	Reaction string
	User     User
}

type TaskProject struct {
	Id   string
	Name string
}

type TaskStatus struct {
	Id     string
	Status string
	Color  string
	Type   string
}

type User struct {
	Id             int
	Username       string
	Color          *string
	Email          string `json:"email"`
	Initials       *string
	ProfilePicture *string
	Role           *int
}

type Folder struct {
	Id     string
	Hidden bool
	Lists  []List
	Name   string
	Space  TaskSpace
}

type List struct {
	Id        string
	Name      string
	Space     TaskSpace
	Status    ListStatus
	DueDate   string `json:"due_date"`
	StartDate string `json:"start_date"`
}

type ListStatus struct {
	Status string
	Color  string
}

type Subtask struct {
	Id        string
	CustomId  *string `json:"custom_id"`
	Name      string
	Status    TaskStatus
	Creator   User
	Assignees []User
	Watchers  []User
}
