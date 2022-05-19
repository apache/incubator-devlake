package main

import (
	"github.com/apache/incubator-devlake/api"
	_ "github.com/apache/incubator-devlake/version"
)

func main() {
	api.CreateApiService()
}
