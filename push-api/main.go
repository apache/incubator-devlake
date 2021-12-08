package main

import (
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

	router.Run("localhost:9123")
}

func Post(c *gin.Context) {
	var err error
	tableName := c.Param("tableName")
	var unknownThings []map[string]interface{}
	err = c.BindJSON(&unknownThings)
	if err != nil {
		c.AbortWithError(400, err)
	}
	for _, unknownThing := range unknownThings {
		err := db.InsertThing(tableName, unknownThing)
		if err != nil {
			c.AbortWithError(500, err)
		}
	}
	c.IndentedJSON(http.StatusOK, []string{"thanks"})
}
