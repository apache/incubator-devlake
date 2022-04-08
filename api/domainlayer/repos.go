package domainlayer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/api/shared"
	"github.com/merico-dev/lake/services"
)

/*
Get all repos from database
GET /repos
{
	"repos": [
		{"id": "github:GithubRepo:384111310", "name": "merico-dev/lake", ...}
	],
	"count": 5
}
*/
func ReposIndex(c *gin.Context) {
	repos, count, err := services.GetRepos()
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	shared.ApiOutputSuccess(c, gin.H{"repos": repos, "count": count}, http.StatusOK)
}
