package main

import (
	"log"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/runner"
	_ "github.com/merico-dev/lake/version"
	"github.com/merico-dev/lake/worker/app"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// basic resources
	cfg := config.GetConfig()
	db, err := runner.NewGormDb(cfg, logger.Global)
	if err != nil {
		panic(err)
	}
	err = runner.LoadPlugins(cfg.GetString("PLUGIN_DIR"), cfg, logger.Global, db)
	if err != nil {
		panic(err)
	}

	// establish temporal connection
	TASK_QUEUE := cfg.GetString("TEMPORAL_TASK_QUEUE")
	// Create the client object just once per process
	c, err := client.NewClient(client.Options{
		HostPort: cfg.GetString("TEMPORAL_URL"),
		Logger:   app.NewTemporalLogger(logger.Global),
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, TASK_QUEUE, worker.Options{})
	w.RegisterWorkflow(app.DevLakePipelineWorkflow)
	w.RegisterActivity(app.DevLakeTaskActivity)
	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
