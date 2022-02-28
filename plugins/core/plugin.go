package core

import (
	"context"
	"net/url"
)

type ApiResourceInput struct {
	Params map[string]string      // path variables
	Query  url.Values             // query string
	Body   map[string]interface{} // json body
}

type ApiResourceOutput struct {
	Body   interface{} // response body
	Status int
}

type ApiResourceHandler func(input *ApiResourceInput) (*ApiResourceOutput, error)

type Plugin interface {
	Description() string
	Init()
	Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error
	// PkgPath information lost when compiled as plugin(.so)
	RootPkgPath() string
	ApiResources() map[string]map[string]ApiResourceHandler
}

type SubTask interface {
	Execute() error
}
