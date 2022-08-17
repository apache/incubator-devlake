package pipelines

import (
	"github.com/apache/incubator-devlake/plugins/core"
	"net/http"
)

type CodeTransformationRules struct {
	PrType               string `mapstructure:"prType" json:"prType"`
	PrComponent          string `mapstructure:"prComponent" json:"prComponent"`
	PrBodyClosePattern   string `mapstructure:"prBodyClosePattern" json:"prBodyClosePattern"`
	IssueSeverity        string `mapstructure:"issueSeverity" json:"issueSeverity"`
	IssuePriority        string `mapstructure:"issuePriority" json:"issuePriority"`
	IssueComponent       string `mapstructure:"issueComponent" json:"issueComponent"`
	IssueTypeBug         string `mapstructure:"issueTypeBug" json:"issueTypeBug"`
	IssueTypeIncident    string `mapstructure:"issueTypeIncident" json:"issueTypeIncident"`
	IssueTypeRequirement string `mapstructure:"issueTypeRequirement" json:"issueTypeRequirement"`
}

type TicketTransformationRules struct {
	EpicKeyField               string `json:"epicKeyField"`
	StoryPointField            string `json:"storyPointField"`
	RemotelinkCommitShaPattern string `json:"remotelinkCommitShaPattern"`
	TypeMappings               map[string]struct {
		StandardType string `json:"standardType"`
	} `json:"typeMappings"`
}

// @Summary pipelines plan for jira
// @Description pipelines plan for jira
// @Tags framework/pipelines
// @Accept application/json
// @Param pipeline-plan body JiraPipelinePlan true "json"
// @Router /pipelines/jira/pipeline-plan [post]
func PostJiraPipeline(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	pipeline := &JiraPipelinePlan{}
	return &core.ApiResourceOutput{Body: pipeline, Status: http.StatusOK}, nil
}

type JiraPipelinePlan [][]struct {
	Plugin   string   `json:"plugin"`
	Subtasks []string `json:"subtasks"`
	Options  struct {
		BoardID             int                       `json:"boardId"`
		ConnectionID        int                       `json:"connectionId"`
		TransformationRules TicketTransformationRules `json:"transformationRules"`
	} `json:"options"`
}

// @Summary pipelines plan for gitextractor
// @Description pipelines plan for gitextractor
// @Tags framework/pipelines
// @Accept application/json
// @Param pipeline body GitextractorPipelinePlan true "json"
// @Router /pipelines/gitextractor/pipeline-plan [post]
func PostGitextractorPipeline(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	pipeline := &GitextractorPipelinePlan{}
	return &core.ApiResourceOutput{Body: pipeline, Status: http.StatusOK}, nil
}

type GitextractorPipelinePlan [][]struct {
	Plugin  string `json:"plugin"`
	Options struct {
		URL    string `json:"url"`
		RepoID string `json:"repoId"`
	} `json:"options"`
}

// @Summary pipelines plan for gitee
// @Description pipelines plan for gitee
// @Tags framework/pipelines
// @Accept application/json
// @Param pipeline body GiteePipelinePlan true "json"
// @Router /pipelines/gitee/pipeline-plan [post]
func PostGiteePipeline(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	pipeline := &GiteePipelinePlan{}
	return &core.ApiResourceOutput{Body: pipeline, Status: http.StatusOK}, nil
}

type GiteePipelinePlan [][]struct {
	Plugin   string   `json:"plugin"`
	Subtasks []string `json:"subtasks"`
	Options  struct {
		ConnectionID   int    `json:"connectionId"`
		Owner          string `json:"owner"`
		Repo           string `json:"repo"`
		Since          string
		Transformation CodeTransformationRules `json:"transformation"`
	} `json:"options"`
}

// @Summary pipelines plan for feishu
// @Description pipelines plan for feishu
// @Tags framework/pipelines
// @Accept application/json
// @Param pipeline body FeishuPipelinePlan true "json"
// @Router /pipelines/feishu/pipeline-plan [post]
func PostFeishuPipeline(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	pipeline := &FeishuPipelinePlan{}
	return &core.ApiResourceOutput{Body: pipeline, Status: http.StatusOK}, nil
}

type FeishuPipelinePlan [][]struct {
	Plugin  string   `json:"plugin"`
	Options struct{} `json:"options"`
}

// @Summary pipelines plan for dbt
// @Description pipelines plan for dbt
// @Tags framework/pipelines
// @Accept application/json
// @Param pipeline body DbtPipelinePlan true "json"
// @Router /pipelines/dbt/pipeline-plan [post]
func PostDbtPipeline(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	pipeline := &DbtPipelinePlan{}
	return &core.ApiResourceOutput{Body: pipeline, Status: http.StatusOK}, nil
}

type DbtPipelinePlan [][]struct {
	Plugin  string `json:"plugin"`
	Options struct {
		ProjectPath    string   `json:"projectPath"`
		ProjectName    string   `json:"projectName"`
		ProjectTarget  string   `json:"projectTarget"`
		SelectedModels []string `json:"selectedModels"`
		ProjectVars    struct {
			Demokey1 string `json:"demokey1"`
			Demokey2 string `json:"demokey2"`
		} `json:"projectVars"`
	} `json:"options"`
}

// @Summary pipelines plan for github
// @Description pipelines plan for github
// @Tags framework/pipelines
// @Accept application/json
// @Param pipeline body GithubPipelinePlan true "json"
// @Router /pipelines/github/pipeline-plan [post]
func PostGithubPipeline(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	pipeline := &GithubPipelinePlan{}
	return &core.ApiResourceOutput{Body: pipeline, Status: http.StatusOK}, nil
}

type GithubPipelinePlan [][]struct {
	Plugin   string   `json:"plugin"`
	Subtasks []string `json:"subtasks"`
	Options  struct {
		ConnectionID   int    `json:"connectionId"`
		Owner          string `json:"owner"`
		Repo           string `json:"repo"`
		Since          string
		Transformation CodeTransformationRules `json:"transformation"`
	} `json:"options"`
}

// @Summary pipelines plan for tapd
// @Description pipelines plan for tapd
// @Tags framework/pipelines
// @Accept application/json
// @Param pipeline body TapdPipelinePlan true "json"
// @Router /pipelines/tapd/pipeline-plan [post]
func PostTapdPipelinePlan(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	pipeline := &TapdPipelinePlan{}
	return &core.ApiResourceOutput{Body: pipeline, Status: http.StatusOK}, nil
}

type TapdPipelinePlan [][]struct {
	Plugin   string   `json:"plugin"`
	Subtasks []string `json:"subtasks"`
	Options  struct {
		WorkspaceId uint64   `mapstruct:"workspaceId"`
		CompanyId   uint64   `mapstruct:"companyId"`
		Tasks       []string `mapstruct:"tasks,omitempty"`
		Since       string
	} `json:"options"`
}
