package main

import "github.com/merico-dev/lake/plugins"

// @title Dev Lake API
// @version 1.0

// @host localhost:8080
// @BasePath /
func main() {
	err := plugins.LoadPlugins("./plugins")
	if err != nil {
		panic(err)
	}
	CreateApiService()
	println("Hello, lake")
}
