package blueprints

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

// @Summary blueprints setting for tapd
// @Description blueprint setting for tapd
// @Tags framework/blueprints
// @Accept application/json
// @Param blueprint body TapdBlueprintSetting true "json"
// @Router /blueprints/tapd/blueprint-setting [post]
func PostTapdBlueprintSetting(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	blueprint := &TapdBlueprintSetting{}
	return &core.ApiResourceOutput{Body: blueprint, Status: http.StatusOK}, nil
}

type TapdBlueprintSetting []struct {
	Version     string `json:"version"`
	Connections []struct {
		Plugin       string `json:"plugin"`
		ConnectionID int    `json:"connectionId"`
		Scope        []struct {
			Options struct {
				WorkspaceId uint64   `mapstruct:"workspaceId"`
				CompanyId   uint64   `mapstruct:"companyId"`
				Tasks       []string `mapstruct:"tasks,omitempty"`
				Since       string
			} `json:"options"`
			Entities []string `json:"entities"`
		} `json:"scope"`
	} `json:"connections"`
}

// @Summary blueprints plan for refdiff
// @Description blueprints plan for refdiff
// @Tags framework/blueprints
// @Accept application/json
// @Param blueprint body RefdiffBlueprintPlan true "json"
// @Router /blueprints/refdiff/blueprint-plan [post]
func PostRefdiffBlueprint(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	blueprint := &RefdiffBlueprintPlan{}
	return &core.ApiResourceOutput{Body: blueprint, Status: http.StatusOK}, nil
}

type RefdiffBlueprintPlan [][]struct {
	Plugin  string `json:"plugin"`
	Options struct {
		RepoID string `json:"repoId"`
		Pairs  []struct {
			NewRef string `json:"newRef"`
			OldRef string `json:"oldRef"`
		} `json:"pairs"`
	} `json:"options"`
}

// @Summary blueprints setting for jira
// @Description blueprint setting for jira
// @Tags framework/blueprints
// @Accept application/json
// @Param blueprint-setting body JiraBlueprintSetting true "json"
// @Router /blueprints/jira/blueprint-setting [post]
func PostJiraBlueprintSetting(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	blueprint := &JiraBlueprintSetting{}
	return &core.ApiResourceOutput{Body: blueprint, Status: http.StatusOK}, nil
}

type JiraBlueprintSetting []struct {
	Version     string `json:"version"`
	Connections []struct {
		Plugin       string `json:"plugin"`
		ConnectionID int    `json:"connectionId"`
		Scope        []struct {
			Transformation TicketTransformationRules `json:"transformation"`
			Options        struct {
				BoardId uint64 `json:"boardId"`
				Since   string `json:"since"`
			} `json:"options"`
			Entities []string `json:"entities"`
		} `json:"scope"`
	} `json:"connections"`
}

// @Summary blueprints setting for gitlab
// @Description blueprint setting for gitlab
// @Tags framework/blueprints
// @Accept application/json
// @Param blueprint body GitlabBlueprintSetting true "json"
// @Router /blueprints/gitlab/blueprint-setting [post]
func PostGitlabBluePrint(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	blueprint := &GitlabBlueprintSetting{}
	return &core.ApiResourceOutput{Body: blueprint, Status: http.StatusOK}, nil
}

type GitlabBlueprintSetting []struct {
	Version     string `json:"version"`
	Connections []struct {
		Plugin       string `json:"plugin"`
		ConnectionID int    `json:"connectionId"`
		Scope        []struct {
			Transformation CodeTransformationRules `json:"transformation"`
			Options        struct {
				ProjectId int `json:"projectId"`
				Since     string
			} `json:"options"`
			Entities []string `json:"entities"`
		} `json:"scope"`
	} `json:"connections"`
}

// @Summary blueprints plan for gitextractor
// @Description blueprints plan for gitextractor
// @Tags framework/blueprints
// @Accept application/json
// @Param blueprint body GitextractorBlueprintPlan true "json"
// @Router /blueprints/gitextractor/blueprint-plan [post]
func PostGitextractorBlueprint(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	blueprint := &GitextractorBlueprintPlan{}
	return &core.ApiResourceOutput{Body: blueprint, Status: http.StatusOK}, nil
}

type GitextractorBlueprintPlan [][]struct {
	Plugin  string `json:"plugin"`
	Options struct {
		URL    string `json:"url"`
		RepoID string `json:"repoId"`
	} `json:"options"`
}

// @Summary blueprints plan for feishu
// @Description blueprints plan for feishu
// @Tags framework/blueprints
// @Accept application/json
// @Param blueprint body FeishuBlueprintPlan true "json"
// @Router /blueprints/feishu/blueprint-plan [post]
func PostFeishuBlueprint(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	blueprint := &FeishuBlueprintPlan{}
	return &core.ApiResourceOutput{Body: blueprint, Status: http.StatusOK}, nil
}

type FeishuBlueprintPlan [][]struct {
	Plugin  string   `json:"plugin"`
	Options struct{} `json:"options"`
}

// @Summary blueprints plan for dbt
// @Description blueprints plan for dbt
// @Tags framework/blueprints
// @Accept application/json
// @Param blueprint body DbtBlueprintPlan true "json"
// @Router /blueprints/dbt/blueprint-plan [post]
func PostDbtBlueprint(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	blueprint := &DbtBlueprintPlan{}
	return &core.ApiResourceOutput{Body: blueprint, Status: http.StatusOK}, nil
}

type DbtBlueprintPlan [][]struct {
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

// @Summary blueprints setting for github
// @Description blueprint setting for github
// @Tags framework/blueprints
// @Accept application/json
// @Param blueprint body GithubBlueprintSetting true "json"
// @Router /blueprints/github/blueprint-setting [post]
func PostGithubBluePrint(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	blueprint := &GithubBlueprintSetting{}
	return &core.ApiResourceOutput{Body: blueprint, Status: http.StatusOK}, nil
}

type GithubBlueprintSetting []struct {
	Version     string `json:"version"`
	Connections []struct {
		Plugin       string `json:"plugin"`
		ConnectionID int    `json:"connectionId"`
		Scope        []struct {
			Transformation CodeTransformationRules `json:"transformation"`
			Options        struct {
				Owner string `json:"owner"`
				Repo  string `json:"repo"`
				Since string
			} `json:"options"`
			Entities []string `json:"entities"`
		} `json:"scope"`
	} `json:"connections"`
}

// @Summary blueprints setting for icla
// @Description blueprint setting for icla
// @Tags framework/blueprints
// @Accept application/json
// @Param blueprint body iclaBlueprintSetting true "json"
// @Router /blueprints/icla/blueprint-setting [post]
func _() {}

type iclaBlueprintSetting []struct {
	Version     string `json:"version" example:"1.0.0"`
	Connections []struct {
		Plugin string `json:"plugin" example:"icla"`
	} `json:"connections"`
}
