package api

import (
	"net/http"

	"github.com/apache/incubator-devlake/api/shared"
	"github.com/apache/incubator-devlake/models"
	"github.com/gin-gonic/gin"
)

// TagResponse is the API response for a tag
type TagResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Tag     models.Tag `json:"tag"`
}

// TagsResponse is the API response for multiple tags
type TagsResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Tags    []models.Tag `json:"tags"`
}

// TagRequest is the request body for creating/updating tags
type TagRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color" default:"#3399FF"`
}

// RegisterTagsRoutes registers the routes for tag management
func RegisterTagsRoutes(router *gin.RouterGroup) {
	// Get all tags
	router.GET("/tags", listTags)
	
	// Create a new tag
	router.POST("/tags", createTag)
	
	// Get a specific tag
	router.GET("/tags/:id", getTag)
	
	// Update a tag
	router.PATCH("/tags/:id", updateTag)
	
	// Delete a tag
	router.DELETE("/tags/:id", deleteTag)
	
	// Project tag association endpoints
	router.POST("/projects/:projectId/tags/:tagId", addTagToProject)
	router.DELETE("/projects/:projectId/tags/:tagId", removeTagFromProject)
	router.GET("/projects/:projectId/tags", getProjectTags)
}

// listTags returns all tags
func listTags(c *gin.Context) {
	var tags []models.Tag
	if err := db.Find(&tags).Error; err != nil {
		shared.ApiErrorHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, TagsResponse{
		Success: true,
		Tags:    tags,
	})
}

// createTag creates a new tag
func createTag(c *gin.Context) {
	var req TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ApiErrorHandler(c, err)
		return
	}
	
	tag := models.Tag{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
	}
	
	if err := db.Create(&tag).Error; err != nil {
		shared.ApiErrorHandler(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, TagResponse{
		Success: true,
		Message: "Tag created successfully",
		Tag:     tag,
	})
}

// getTag returns a specific tag by ID
func getTag(c *gin.Context) {
	id := c.Param("id")
	var tag models.Tag
	
	if err := db.First(&tag, "id = ?", id).Error; err != nil {
		shared.ApiErrorHandler(c, err)
		return
	}
	
	c.JSON(http.StatusOK, TagResponse{
		Success: true,
		Tag:     tag,
	})
}

// updateTag updates a tag
func updateTag(c *gin.Context) {
	id := c.Param("id")
	var req TagRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ApiErrorHandler(c, err)
		return
	}
	
	var tag models.Tag
	if err := db.First(&tag, "id = ?", id).Error; err != nil {
		shared.ApiErrorHandler(c, err)
		return
	}
	
	tag.Name = req.Name
	tag.Description = req.Description
	tag.Color = req.Color
	
	if err := db.Save(&tag).Error; err != nil {
		shared.ApiErrorHandler(c, err)
		return
	}
	
	c.JSON(http.StatusOK, TagResponse{
		Success: true,
		Message: "Tag updated successfully",
		Tag:     tag,
	})
}

// deleteTag deletes a tag
func deleteTag(c *gin.Context) {
	id := c.Param("id")
	
	// Delete tag associations first
	if err := db.Delete(&models.ProjectTag{}, "tag_id = ?", id).Error; err != nil {
		shared.ApiErrorHandler(c, err)
		return
	}
	
	// Delete the tag
	if err := db.Delete(&models.Tag{}, "id = ?", id).Error; err != nil {
		shared.ApiErrorHandler(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tag deleted successfully",
	})
}

// addTagToProject associates a tag with a project
func addTagToProject(c *gin.Context) {
	projectId := c.Param("projectId")
	tagId := c.Param("tagId")
	
	// Check if project exists
	var project models.Project
	if err := db.First(&project, "id = ?", projectId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	
	// Check if tag exists
	var tag models.Tag
	if err := db.First(&tag, "id = ?", tagId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}
	
	// Create association
	projectTag := models.ProjectTag{
		ProjectId: projectId,
		TagId:     tagId,
	}
	
	// Check if association already exists
	var count int64
	db.Model(&models.ProjectTag{}).Where("project_id = ? AND tag_id = ?", projectId, tagId).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Project already has this tag"})
		return
	}
	
	if err := db.Create(&projectTag).Error; err != nil {
		shared.ApiErrorHandler(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tag added to project successfully",
	})
}

// removeTagFromProject removes a tag from a project
func removeTagFromProject(c *gin.Context) {
	projectId := c.Param("projectId")
	tagId := c.Param("tagId")
	
	if err := db.Delete(&models.ProjectTag{}, "project_id = ? AND tag_id = ?", projectId, tagId).Error; err != nil {
		shared.ApiErrorHandler(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tag removed from project successfully",
	})
}

// getProjectTags gets all tags for a project
func getProjectTags(c *gin.Context) {
	projectId := c.Param("projectId")
	
	// Check if project exists
	var project models.Project
	if err := db.First(&project, "id = ?", projectId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	
	var tags []models.Tag
	if err := db.Model(&project).Association("Tags").Find(&tags); err != nil {
		shared.ApiErrorHandler(c, err)
		return
	}
	
	c.JSON(http.StatusOK, TagsResponse{
		Success: true,
		Tags:    tags,
	})
}
