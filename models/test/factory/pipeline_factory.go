package factory

import (
	"github.com/merico-dev/lake/models"
)

func CreatePipeline() (*models.Pipeline, error) {
	pipeline := &models.Pipeline{
		Name:          "My Pipeline",
		FinishedTasks: RandInt(),
		Status:        "MY_STATUS",
		Message:       "",
		SpentSeconds:  RandInt(),
	}
	return pipeline, nil
}
