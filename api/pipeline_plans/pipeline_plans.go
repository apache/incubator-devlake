package pipelineplans

import (
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
	inputPipelinePlan := &models.InputPipelinePlan{}

	err := ctx.MustBindWith(inputPipelinePlan, binding.JSON)
	if err != nil {
		logger.Error("post /pipeline failed", err)
		shared.ApiOutputError(ctx, err, http.StatusBadRequest)
		return
	}

	pipelinePlan, err := services.CreatePipelinePlan(inputPipelinePlan)
	if err != nil {
		shared.ApiOutputError(ctx, err, http.StatusBadRequest)
		return
	}

	shared.ApiOutputSuccess(ctx, pipelinePlan, http.StatusCreated)
}

func Index(ctx *gin.Context) {
	pipelinePlans, count, err := services.GetPipelinePlans()
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	shared.ApiOutputSuccess(ctx, gin.H{"pipeline-plans": pipelinePlans, "count": count}, http.StatusOK)
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

func PUT(ctx *gin.Context) {
	pipelinePlanId := ctx.Param("pipelinePlanId")
	id, err := strconv.ParseUint(pipelinePlanId, 10, 64)
	if err != nil {
		shared.ApiOutputError(ctx, err, http.StatusBadRequest)
		return
	}
	editPipelinePlan := &models.EditPipelinePlan{}
	err = ctx.MustBindWith(editPipelinePlan, binding.JSON)
	if err != nil {
		logger.Error("patch /pipeline-plans/:pipelinePlanId failed", err)
		shared.ApiOutputError(ctx, err, http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	editPipelinePlan.PipelinePlanId = id
	pipelinePlan, err := services.ModifyPipelinePlan(editPipelinePlan)
	if err != nil {
		shared.ApiOutputError(ctx, err, http.StatusBadRequest)
		return
	}
	shared.ApiOutputSuccess(ctx, pipelinePlan, http.StatusOK)
}
