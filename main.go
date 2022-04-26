package main

import (
	"github.com/merico-dev/lake/api"
	_ "github.com/merico-dev/lake/version"
)

func main() {
	api.CreateApiService()
}
