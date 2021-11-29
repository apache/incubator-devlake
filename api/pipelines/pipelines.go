package pipelines

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/merico-dev/lake/logger"
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

func Post(ctx *gin.Context) {
	newPipeline := &models.NewPipeline{}

	err := ctx.MustBindWith(newPipeline, binding.JSON)
	if err != nil {
		logger.Error("post /pipeline failed", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	pipeline, err := services.CreatePipeline(newPipeline)
	// Return all created tasks to the User
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	go func() {
		_ = services.RunPipeline(pipeline.ID)
	}()
	ctx.JSON(http.StatusCreated, pipeline)
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
func Index(ctx *gin.Context) {
	var query services.PipelineQuery
	err := ctx.BindQuery(&query)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	pipelines, count, err := services.GetPipelines(&query)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"pipelines": pipelines, "count": count})
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
func Get(ctx *gin.Context) {
	pipelineId := ctx.Param("pipelineId")
	id, err := strconv.ParseUint(pipelineId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid id")
		return
	}
	pipeline, err := services.GetPipeline(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, pipeline)
}

/*
Cancel a pending pipeline
DELETE /pipelines/:pipelineId
*/
func Delete(ctx *gin.Context) {
	pipelineId := ctx.Param("pipelineId")
	id, err := strconv.ParseUint(pipelineId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid id")
		return
	}
	err = services.CancelPipeline(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	}
}
