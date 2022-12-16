/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"log"

	"github.com/apache/incubator-devlake/errors"

	"github.com/apache/incubator-devlake/logger"
	"github.com/apache/incubator-devlake/runner"
	_ "github.com/apache/incubator-devlake/version"
	"github.com/apache/incubator-devlake/worker/app"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	basicRes := runner.CreateAppBasicRes()
	err := runner.LoadPlugins(basicRes)
	if err != nil {
		panic(err)
	}

	// establish temporal connection
	TASK_QUEUE := basicRes.GetConfig("TEMPORAL_TASK_QUEUE")
	// Create the client object just once per process
	c, err := errors.Convert01(client.NewClient(client.Options{
		HostPort: basicRes.GetConfig("TEMPORAL_URL"),
		Logger:   app.NewTemporalLogger(logger.Global),
	}))
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, TASK_QUEUE, worker.Options{})
	w.RegisterWorkflow(app.DevLakePipelineWorkflow)
	w.RegisterActivity(app.DevLakeTaskActivity)
	// Start listening to the Task Queue
	err = errors.Convert(w.Run(worker.InterruptCh()))
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
