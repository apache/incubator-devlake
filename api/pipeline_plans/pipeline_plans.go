package pipelineplans

import (
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/api/shared"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/services"
)

func Post(ctx *gin.Context) {
	newPipeline := &models.NewPipeline{}

	err := ctx.MustBindWith(newPipeline, binding.JSON)
	if err != nil {
		logger.Error("post /pipeline failed", err)
		shared.ApiOutputError(ctx, err, http.StatusBadRequest)
		return
	}

	pipelinePlan, err := services.CreatePipelinePlan(newPipeline)
	if err != nil {
		shared.ApiOutputError(ctx, err, http.StatusBadRequest)
		return
	}

	var tasks [][]*models.NewTask
	err = json.Unmarshal(pipelinePlan.Tasks, &tasks)
	if err != nil {
		shared.ApiOutputError(ctx, err, http.StatusBadRequest)
		return
	}

	pipeline, err := services.CreatePipeline(pipelinePlan.Name, tasks, pipelinePlan.ID)
	// Return all created tasks to the User
	if err != nil {
		shared.ApiOutputError(ctx, err, http.StatusBadRequest)
		return
	}

	go func() {
		_ = services.RunPipeline(pipeline.ID)
	}()
	shared.ApiOutputSuccess(ctx, pipelinePlan, http.StatusCreated)
}

func Index(ctx *gin.Context) {
	pipelines, count, err := services.GetPipelinePlans()
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	shared.ApiOutputSuccess(ctx, gin.H{"pipelines": pipelines, "count": count}, http.StatusOK)
}

func Get(ctx *gin.Context) {
	pipelinePlanId := ctx.Param("pipelinePlanId")
	id, err := strconv.ParseUint(pipelinePlanId, 10, 64)
	if err != nil {
		shared.ApiOutputError(ctx, err, http.StatusBadRequest)
		return
	}
	pipelinePlan, err := services.GetPipelinePlan(id)
	if err != nil {
		shared.ApiOutputError(ctx, err, http.StatusBadRequest)
		return
	}
	shared.ApiOutputSuccess(ctx, pipelinePlan, http.StatusOK)
}

func Delete(ctx *gin.Context) {
	pipelineId := ctx.Param("pipelinePlanId")
	id, err := strconv.ParseUint(pipelineId, 10, 64)
	if err != nil {
		shared.ApiOutputError(ctx, err, http.StatusBadRequest)
		return
	}
	err = services.DeletePipelinePlan(id)
	if err != nil {
		shared.ApiOutputError(ctx, err, http.StatusBadRequest)
	}
}

func Patch(ctx *gin.Context) {
	pipelinePlanId := ctx.Param("pipelinePlanId")
	id, err := strconv.ParseUint(pipelinePlanId, 10, 64)
	if err != nil {
		shared.ApiOutputError(ctx, err, http.StatusBadRequest)
		return
	}
	newPipeline := &models.NewPipeline{}
	err = ctx.MustBindWith(newPipeline, binding.JSON)
	if err != nil {
		logger.Error("patch /pipelines/plans/:pipelinePlanId failed", err)
		shared.ApiOutputError(ctx, err, http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	pipelinePlan, err := services.ModifyPipelinePlan(newPipeline, id)
	if err != nil {
		shared.ApiOutputError(ctx, err, http.StatusBadRequest)
		return
	}
	shared.ApiOutputSuccess(ctx, pipelinePlan, http.StatusOK)
}
