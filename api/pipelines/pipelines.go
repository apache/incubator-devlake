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
// @Summary Create and run a new pipeline
// @Description Create and run a new pipeline
// @Tags pipelines
// @Accept application/json
// @Param pipeline body string true "pipline is json format"
// @Success 200  {object} models.Pipeline
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /pipelines [post]
func Post(c *gin.Context) {
	newPipeline := &models.NewPipeline{}

	err := c.MustBindWith(newPipeline, binding.JSON)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}

	pipeline, err := services.CreatePipeline(newPipeline)
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
GET /pipelines?status=TASK_RUNNING&pending=1&page=1&pagesize=10
{
	"pipelines": [
		{"id": 1, "name": "test-pipeline", ...}
	],
	"count": 5
}
*/

// @Summary Get list of pipelines
// @Description GET /pipelines?status=TASK_RUNNING&pending=1&page=1&pagesize=10
// @Tags pipelines
// @Param status query string true "query"
// @Param pending query int true "query"
// @Param page query int true "query"
// @Param pagesize query int true "query"
// @Success 200  {object} models.Pipeline
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /pipelines [get]
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
// @Get detail of a pipeline
// @Description GET /pipelines/:pipelineId
// @Tags pipelines
// @Param pipelineId path int true "query"
// @Success 200  {object} models.Pipeline
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /pipelines/{pipelineId} [get]
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
// @Cancel a pending pipeline
// @Description DELETE /pipelines/:pipelineId
// @Tags pipelines
// @Param pipelineId path int true "id of your pipline"
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /pipelines/{pipelineId} [delete]
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
