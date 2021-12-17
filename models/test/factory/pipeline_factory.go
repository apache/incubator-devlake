package factory

import (
	"github.com/merico-dev/lake/models"
)

func CreatePipeline() (*models.Pipeline, error) {
	pipeline := &models.Pipeline{
		Name:          "My Pipeline",
		FinishedTasks: 0,
		Status:        "MY_STATUS",
		Message:       "",
		SpentSeconds:  0,
	}
	return pipeline, nil
}
