package pipelineplans

import (
	"fmt"
	"io/ioutil"
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

	var pipelinePlan *models.PipelinePlan
	if newPipeline.CronConfig != nil {
		pipelinePlan, err = services.CreatePipelinePlan(newPipeline)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
	}

	pipeline, err := services.CreatePipeline(newPipeline, pipelinePlan)
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
GET /pipelinePlans
{
	"pipelines": [
		{"id": 1, "name": "test-pipeline", ...}
	],
	"count": 5
}
*/
func Index(ctx *gin.Context) {
	pipelines, count, err := services.GetPipelinePlans()
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
	pipelinePlanId := ctx.Param("pipelinePlanId")
	id, err := strconv.ParseUint(pipelinePlanId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid id")
		return
	}
	pipelinePlan, err := services.GetPipelinePlan(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, pipelinePlan)
}

/*
Cancel a pending pipeline
DELETE /pipelines/:pipelineId
*/
func Delete(ctx *gin.Context) {
	pipelineId := ctx.Param("pipelinePlanId")
	id, err := strconv.ParseUint(pipelineId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid id")
		return
	}
	err = services.DeletePipelinePlan(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	}
}

func Patch(ctx *gin.Context) {
	pipelinePlanId := ctx.Param("pipelinePlanId")
	id, err := strconv.ParseUint(pipelinePlanId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid id")
		return
	}
	newPipeline := &models.NewPipeline{}
	r, _ := ioutil.ReadAll(ctx.Request.Body)

	fmt.Println(string(r))
	err = ctx.MustBindWith(newPipeline, binding.JSON)
	if err != nil {
		logger.Error("patch /pipelines/plans/:pipelinePlanId failed", err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	newPipeline.PipelinePlanId = id
	pipelinePlan, err := services.ModifyPipelinePlan(newPipeline)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, pipelinePlan)
}
