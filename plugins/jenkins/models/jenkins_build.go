package models

import "github.com/merico-dev/lake/models"

type JenkinsBuild struct {
	models.Model
	Duration          float64
	DisplayName       string
	EstimatedDuration float64
	Number            int64
	Result            string
	Timestamp         int64
}
