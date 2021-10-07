package utils

import (
	"fmt"

	"github.com/merico-dev/lake/logger"
	"github.com/nats-io/nats.go"
)

func ListenForCancelEvent(pluginName string, scheduler *WorkerScheduler, progress chan<- float32) {
	// Simple Async Subscriber for cancelling collection
	nc, _ := nats.Connect(nats.DefaultURL)
	_, errNc := nc.Subscribe(pluginName, func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
		scheduler.Release()
		progress <- 1
		logger.Info("You cancelled the collector with nats", false)
		close(progress)
	})

	if errNc != nil {
		logger.Error("errNc", errNc)
	}
}
