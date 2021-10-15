package models

import "github.com/merico-dev/lake/models"

type Pokemon struct {
	Name           string
	URL            string
	BaseExperience int
	models.NoPKModel
}
