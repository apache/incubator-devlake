package main

import (
	"strings"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"myapp/models"
	"myapp/shared"
)

var db *gorm.DB

// listProjects returns all projects with optional tag filtering
func listProjects(c *gin.Context) {
    // Get tag filter params if any
    tagFilter := c.Query("tags")
    
    var projects []models.Project
    query := db.Model(&models.Project{})
    
    // Apply tag filtering if provided
    if tagFilter != "" {
        tags := strings.Split(tagFilter, ",")
        query = query.Joins("JOIN _devlake_project_tags pt ON pt.project_id = projects.id").
              Joins("JOIN tags t ON t.id = pt.tag_id").
              Where("t.name IN ?", tags).
              Group("projects.id").
              Having("COUNT(DISTINCT t.name) = ?", len(tags))
    }
    
    // Execute the query
    if err := query.Find(&projects).Error; err != nil {
        shared.ApiErrorHandler(c, err)
        return
    }
    
    // Load the tags for each project
    for i := range projects {
        db.Model(&projects[i]).Association("Tags").Find(&projects[i].Tags)
    }
    
    c.JSON(200, projects)
}

// getProject returns a specific project
func getProject(c *gin.Context) {
    var project models.Project
    if err := db.First(&project, c.Param("id")).Error; err != nil {
        shared.ApiErrorHandler(c, err)
        return
    }
    
    // Load associated tags
    db.Model(&project).Association("Tags").Find(&project.Tags)
    
    c.JSON(200, project)
}

func main() {
    r := gin.Default()
    r.GET("/projects", listProjects)
    r.GET("/projects/:id", getProject)
    r.Run()
}
