package factory

import (
	"github.com/merico-dev/lake/models"
)

func CreateTask(pipelineId uint64) (*models.Task, error) {
	task := &models.Task{
		Plugin:     "myPlugin",
		Options:    nil,
		Status:     "MY_STATUS",
		Message:    "my message",
		PipelineId: pipelineId,
	}
	return task, nil
}
