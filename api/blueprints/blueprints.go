package blueprints

import (
	"net/http"
	"strconv"

	"github.com/merico-dev/lake/api/shared"

	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/services"
)

// @Summary Post BluePrints
// @Description Post BluePrints
// @Tags blueprints
// @Accept application/json
// @Param blueprint body string true "json"
// @Success 200  {object} models.Blueprint
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /blueprints [post]
func Post(c *gin.Context) {
	blueprint := &models.Blueprint{}

	err := c.ShouldBind(blueprint)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}

	err = services.CreateBlueprint(blueprint)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}

	shared.ApiOutputSuccess(c, blueprint, http.StatusCreated)
}

// @Summary Get BluePrints
// @Description Get BluePrints
// @Tags blueprints
// @Accept application/json
// @Success 200  {object} models.Blueprint
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /blueprints [get]
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

// @Summary get BluePrints
// @Description get BluePrints
// @Tags blueprints
// @Accept application/json
// @Param blueprintId path int true "blueprint id"
// @Success 200  {object} models.Blueprint
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /blueprints/{blueprintId} [get]
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

// @Summary Delete BluePrints
// @Description Delete BluePrints
// @Tags blueprints
// @Param blueprintId path string true "blueprintId"
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /blueprints/{blueprintId} [delete]
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

/*
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
*/

// @Summary Patch BluePrints
// @Description Patch BluePrints
// @Tags blueprints
// @Accept application/json
// @Param blueprintId path string true "blueprintId"
// @Success 200  {object} models.Blueprint
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /blueprints/{blueprintId} [Patch]
func Patch(c *gin.Context) {
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
	err = c.ShouldBind(blueprint)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	err = services.UpdateBlueprint(blueprint)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	shared.ApiOutputSuccess(c, blueprint, http.StatusOK)
}
