package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jakkayy/envSync/internal/database"
)

type CreateProjectRequest struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	proj := database.Project{
		ID:            req.ID,
		Name:          req.Name,
		SecretKeyHash: "secret-hash-placeholder",
	}

	if err := database.DB.Create(&proj).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "project already exists or database error"})
		return
	}

	defaultEnvs := []string{"dev", "staging", "prod"}
	for _, envName := range defaultEnvs {
		database.DB.Create(&database.Environment{
			ProjectID:     proj.ID,
			Name:          envName,
			LatestVersion: 0,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "project created successfully",
		"project": proj,
	})
}
