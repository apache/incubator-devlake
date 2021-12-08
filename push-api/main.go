package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/push-api/db"
)

// Next steps
// Connect to DB [x]
// Insert record (commit) into DB using SQL

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
	tableName := c.Param("tableName")
	var unknownThings []map[string]interface{}
	err = c.BindJSON(&unknownThings)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
	}
	for _, unknownThing := range unknownThings {
		err := db.InsertThing(tableName, unknownThing)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
	}
	c.IndentedJSON(http.StatusOK, []string{"thanks"})
}
