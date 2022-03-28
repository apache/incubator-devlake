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

func Post(c *gin.Context) {
	inputPipelinePlan := &models.InputPipelinePlan{}

	err := c.MustBindWith(inputPipelinePlan, binding.JSON)
	if err != nil {
		logger.Error("post /pipeline failed", err)
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}

	pipelinePlan, err := services.CreatePipelinePlan(inputPipelinePlan)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}

	shared.ApiOutputSuccess(c, pipelinePlan, http.StatusCreated)
}

func Index(c *gin.Context) {
	var query services.PipelinePlanQuery
	err := c.ShouldBindQuery(&query)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	pipelinePlans, count, err := services.GetPipelinePlans(&query)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	shared.ApiOutputSuccess(c, gin.H{"pipeline-plans": pipelinePlans, "count": count}, http.StatusOK)
}

func Get(c *gin.Context) {
	pipelinePlanId := c.Param("pipelinePlanId")
	id, err := strconv.ParseUint(pipelinePlanId, 10, 64)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	pipelinePlan, err := services.GetPipelinePlan(id)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	shared.ApiOutputSuccess(c, pipelinePlan, http.StatusOK)
}

func Delete(c *gin.Context) {
	pipelineId := c.Param("pipelinePlanId")
	id, err := strconv.ParseUint(pipelineId, 10, 64)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	err = services.DeletePipelinePlan(id)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
	}
}

func PUT(c *gin.Context) {
	pipelinePlanId := c.Param("pipelinePlanId")
	id, err := strconv.ParseUint(pipelinePlanId, 10, 64)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	editPipelinePlan := &models.EditPipelinePlan{}
	err = c.MustBindWith(editPipelinePlan, binding.JSON)
	if err != nil {
		logger.Error("patch /pipeline-plans/:pipelinePlanId failed", err)
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	editPipelinePlan.PipelinePlanId = id
	pipelinePlan, err := services.ModifyPipelinePlan(editPipelinePlan)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	shared.ApiOutputSuccess(c, pipelinePlan, http.StatusOK)
}
