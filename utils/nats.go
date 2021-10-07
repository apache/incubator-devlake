package utils

import (
	"fmt"

	"github.com/merico-dev/lake/logger"
	"github.com/nats-io/nats.go"
)

func ListenForCancelEvent(scheduler *WorkerScheduler, progress chan<- float32, taskId uint64) {
	// Simple Async Subscriber for cancelling collection
	nc, _ := nats.Connect(nats.DefaultURL)
	_, errNc := nc.Subscribe("cancelTask", func(message *nats.Msg) {
		// fmt.Printf("INFO: Received a message: %v\n", string(message.Data))
		// fmt.Printf("INFO: taskId: %v\n", fmt.Sprint(taskId))
		// fmt.Printf("INFO: taskId: %v\n", string(message.Data) == fmt.Sprint(taskId))

		if string(message.Data) == fmt.Sprint(taskId) {
			scheduler.Release()
			progress <- 1
			logger.Info("INFO: You cancelled task with ID: ", taskId)
			close(progress)
		}
	})

	if errNc != nil {
		logger.Error("errNc", errNc)
	}
}
