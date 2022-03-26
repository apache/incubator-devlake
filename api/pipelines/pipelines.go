package pipelines

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/merico-dev/lake/api/shared"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/services"
)

/*
Create and run a new pipeline
POST /pipelines
{
	"name": "name-of-pipeline",
	"tasks": [
		[ {"plugin": "gitlab", ...}, {"plugin": "jira"} ],
		[ {"plugin": "github", ...}],
	]
}
*/

func Post(c *gin.Context) {
	newPipeline := &models.NewPipeline{}

	err := c.MustBindWith(newPipeline, binding.JSON)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}

	pipeline, err := services.CreatePipeline(newPipeline, nil)
	// Return all created tasks to the User
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}

	go func() {
		_ = services.RunPipeline(pipeline.ID)
	}()
	shared.ApiOutputSuccess(c, pipeline, http.StatusCreated)
}

/*
Get list of pipelines
GET /pipelines?status=TASK_RUNNING&pending=1&page=1&=pagesize=10
{
	"pipelines": [
		{"id": 1, "name": "test-pipeline", ...}
	],
	"count": 5
}
*/
func Index(c *gin.Context) {
	var query services.PipelineQuery
	err := c.ShouldBindQuery(&query)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	pipelines, count, err := services.GetPipelines(&query)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	shared.ApiOutputSuccess(c, gin.H{"pipelines": pipelines, "count": count}, http.StatusOK)
}

/*
Get detail of a pipeline
GET /pipelines/:pipelineId
{
	"id": 1,
	"name": "test-pipeline",
	...
}
*/
func Get(c *gin.Context) {
	pipelineId := c.Param("pipelineId")
	id, err := strconv.ParseUint(pipelineId, 10, 64)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	pipeline, err := services.GetPipeline(id)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	shared.ApiOutputSuccess(c, pipeline, http.StatusOK)
}

/*
Cancel a pending pipeline
DELETE /pipelines/:pipelineId
*/
func Delete(c *gin.Context) {
	pipelineId := c.Param("pipelineId")
	id, err := strconv.ParseUint(pipelineId, 10, 64)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	err = services.CancelPipeline(id)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	shared.ApiOutputSuccess(c, nil, http.StatusOK)
}
