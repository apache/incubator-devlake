package core

import "net/url"

// Contains api request information
type ApiResourceInput struct {
	Params map[string]string      // path variables
	Query  url.Values             // query string
	Body   map[string]interface{} // json body
}

// Describe response data of a api
type ApiResourceOutput struct {
	Body   interface{} // response body
	Status int
}

type ApiResourceHandler func(input *ApiResourceInput) (*ApiResourceOutput, error)

// Implement this interface if plugin offered API
// Code sample to register a api on `sources/:sourceId`:
// func (plugin Jira) ApiResources() map[string]map[string]core.ApiResourceHandler {
// 	return map[string]map[string]core.ApiResourceHandler{
// 		"sources/:sourceId": {
// 			"PUT":    api.PutSource,
// 			"DELETE": api.DeleteSource,
// 			"GET":    api.GetSource,
// 		},
// 	}
// }
type PluginApi interface {
	ApiResources() map[string]map[string]ApiResourceHandler
}
