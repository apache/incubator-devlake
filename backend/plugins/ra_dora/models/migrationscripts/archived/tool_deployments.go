package archived

import "github.com/apache/incubator-devlake/core/models/common"

// Deployment represents the entire JSON structure
type Deployment struct {
	common.NoPKModel `swaggerignore:"true" json:"-" mapstructure:"-"`

	Metadata Metadata `json:"metadata"`
	Spec     Spec     `json:"spec"`
	Status   Status   `json:"status"`
}

// Metadata represents the metadata part of the JSON
type Metadata struct {
	Name              string            `json:"name"`
	GenerateName      string            `json:"generateName"`
	Namespace         string            `json:"namespace"`
	UID               string            `json:"uid"`
	ResourceVersion   string            `json:"resourceVersion"`
	Generation        int               `json:"generation"`
	CreationTimestamp string            `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
	ManagedFields     []ManagedField    `json:"managedFields"`
}

// ManagedField represents the managed fields part of the JSON
type ManagedField struct {
	Manager    string   `json:"manager"`
	Operation  string   `json:"operation"`
	APIVersion string   `json:"apiVersion"`
	Time       string   `json:"time"`
	FieldsType string   `json:"fieldsType"`
	FieldsV1   FieldsV1 `json:"fieldsV1"`
}

// FieldsV1 represents the fieldsV1 part of the JSON
type FieldsV1 struct {
	Metadata MetadataFields `json:"f:metadata"`
	Spec     interface{}    `json:"f:spec"`
	Status   interface{}    `json:"f:status"`
}

// MetadataFields represents the metadata fields part of the JSON
type MetadataFields struct {
	GenerateName interface{} `json:"f:generateName"`
	Labels       Labels      `json:"f:labels"`
	Annotations  Annotations `json:"f:annotations"`
}

// Labels represents the labels part of the JSON
type Labels struct {
	EventsActionTimestamp interface{} `json:"f:events.argoproj.io/action-timestamp"`
	EventsSensor          interface{} `json:"f:events.argoproj.io/sensor"`
	EventsTrigger         interface{} `json:"f:events.argoproj.io/trigger"`
	WorkflowsCreator      interface{} `json:"f:workflows.argoproj.io/creator"`
	WorkflowsCompleted    interface{} `json:"f:workflows.argoproj.io/completed"`
	WorkflowsPhase        interface{} `json:"f:workflows.argoproj.io/phase"`
}

// Annotations represents the annotations part of the JSON
type Annotations struct {
	PodNameFormat interface{} `json:"f:workflows.argoproj.io/pod-name-format"`
}

// Spec represents the spec part of the JSON
type Spec struct {
	Entrypoint          string              `json:"entrypoint"`
	Arguments           Arguments           `json:"arguments"`
	WorkflowTemplateRef WorkflowTemplateRef `json:"workflowTemplateRef"`
}

// Arguments represents the arguments part of the JSON
type Arguments struct {
	Parameters []Parameter `json:"parameters"`
}

// Parameter represents a single parameter in the arguments
type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// WorkflowTemplateRef represents the workflow template reference part of the JSON
type WorkflowTemplateRef struct {
	Name string `json:"name"`
}

// Status represents the status part of the JSON
type Status struct {
	Phase      string                 `json:"phase"`
	StartedAt  string                 `json:"startedAt"`
	FinishedAt string                 `json:"finishedAt"`
	Progress   string                 `json:"progress"`
	Nodes      map[string]interface{} `json:"nodes"`
}

func (Deployment) TableName() string {
	return "_tool_argo_api_deployments"
}
