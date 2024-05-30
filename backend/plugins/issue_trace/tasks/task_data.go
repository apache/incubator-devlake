package tasks

// Options original parameter from bp (or pipeline)
type Options struct {
	Plugin       string `json:"plugin"`       // jira
	ConnectionId uint64 `json:"connectionId"` // 1
	BoardId      uint64 `json:"boardId"`      // 68
	LakeBoardId  string `json:"lakeBoardId"`  // jira:JiraBoard:1:68
}

// TaskData converted parameter
type TaskData struct {
	Options Options
	BoardId string // jira:1:JiraBoard:68
}
