package common

import "net/http"

type ApiAsyncCallback func(*http.Response) error
