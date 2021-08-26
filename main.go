package main

import "github.com/merico-dev/lake/config"

func main() {
	config.ReadConfig()
	CreateApiService()
	println("Hello, lake")
}
