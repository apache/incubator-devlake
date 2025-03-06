package tasks

import (
	"github.com/apache/incubator-devlake/core/errors"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type ArgoOptions struct {
	ConnectionId uint64 `json:"connectionId" mapstructure:"connectionId,omitempty"`
	FullName     string `json:"string" mapstructure:"string,omitempty"`
}

type ArgoTaskData struct {
	Options   *ArgoOptions
	ApiClient *helper.ApiClient
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*ArgoOptions, errors.Error) {
	op, err := DecodeTaskOptions(options)
	if err != nil {
		return nil, err
	}
	err = ValidateTaskOptions(op)
	if err != nil {
		return nil, err
	}
	return op, nil
}

func DecodeTaskOptions(options map[string]interface{}) (*ArgoOptions, errors.Error) {
	var op ArgoOptions
	err := helper.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func ValidateTaskOptions(op *ArgoOptions) errors.Error {
	if op.ConnectionId == 0 {
		return errors.BadInput.New("connectionId is invalid")
	}
	return nil
}
