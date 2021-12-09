package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/push-api/db"
)

/*
	POST /api/:tableName
	[
		{
			"id": 1,
			"sha": "osidjfoawehfwh08"
		}
	]
*/

func main() {
	db.InitDb()
	router := gin.Default()
	router.POST("/api/:tableName", Post)

	err := router.Run("localhost:9123")
	if err != nil {
		fmt.Println("ERROR >>> a2kj4: ", err)
	}
}

func Post(c *gin.Context) {
	var err error
	var totalRowsAffected int64
	tableName := c.Param("tableName")
	var unknownThings []map[string]interface{}
	err = c.BindJSON(&unknownThings)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
	}
	for _, unknownThing := range unknownThings {
		rowsAffected, err := db.InsertThing(tableName, unknownThing)
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
