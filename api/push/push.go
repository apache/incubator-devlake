package push

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/services"
)

/*
	POST /push/:tableName
	[
		{
			"id": 1,
			"sha": "osidjfoawehfwh08"
		}
	]
*/

func Post(c *gin.Context) {
	var err error
	var totalRowsAffected int64
	tableName := c.Param("tableName")
	var rowsToInsert []map[string]interface{}
	err = c.BindJSON(&rowsToInsert)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
	}
	for _, rowToInsert := range rowsToInsert {
		rowsAffected, err := services.InsertRow(tableName, rowToInsert)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
		totalRowsAffected += rowsAffected
	}
	if len(c.Errors) > 0 {
		c.JSON(http.StatusOK, c.Errors)
	} else {
		c.JSON(http.StatusOK, gin.H{"Rows affected": totalRowsAffected})
	}
}
