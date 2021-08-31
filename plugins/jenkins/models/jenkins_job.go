package models

import "github.com/merico-dev/lake/models"

type JenkinsJob struct {
	models.Model
	Name  string
	Class string
	Color string
}
