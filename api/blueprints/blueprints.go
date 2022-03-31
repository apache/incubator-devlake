package blueprints

import (
	"fmt"
	"github.com/merico-dev/lake/api/shared"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/services"
)

func Post(c *gin.Context) {
	inputBlueprint := &models.InputBlueprint{}

	err := c.MustBindWith(inputBlueprint, binding.JSON)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}

	blueprint, err := services.CreateBlueprint(inputBlueprint)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}

	shared.ApiOutputSuccess(c, blueprint, http.StatusCreated)
}

func Index(c *gin.Context) {
	var query services.BlueprintQuery
	err := c.ShouldBindQuery(&query)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	blueprints, count, err := services.GetBlueprints(&query)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	shared.ApiOutputSuccess(c, gin.H{"blueprints": blueprints, "count": count}, http.StatusOK)
}

func Get(c *gin.Context) {
	blueprintId := c.Param("blueprintId")
	id, err := strconv.ParseUint(blueprintId, 10, 64)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	blueprint, err := services.GetBlueprint(id)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	shared.ApiOutputSuccess(c, blueprint, http.StatusOK)
}

func Delete(c *gin.Context) {
	pipelineId := c.Param("blueprintId")
	id, err := strconv.ParseUint(pipelineId, 10, 64)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	err = services.DeleteBlueprint(id)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
	}
}

func Put(c *gin.Context) {
	blueprintId := c.Param("blueprintId")
	id, err := strconv.ParseUint(blueprintId, 10, 64)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	editBlueprint := &models.EditBlueprint{}
	err = c.MustBindWith(editBlueprint, binding.JSON)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	editBlueprint.BlueprintId = id
	blueprint, err := services.ModifyBlueprint(editBlueprint)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	shared.ApiOutputSuccess(c, blueprint, http.StatusOK)
}
