package main

import "github.com/merico-dev/lake/plugins"

func main() {
	err := plugins.LoadPlugins("./plugins")
	if err != nil {
		panic(err)
	}
	CreateApiService()
	println("Hello, lake")
}
