package tasks

type AeOptions struct {
	ProjectId int
	Tasks     []string `json:"tasks,omitempty"`
}

type AeTaskData struct {
	Options   *AeOptions
	ApiClient *AEApiClient
}
type AeApiParams struct {
	ProjectId int
}
