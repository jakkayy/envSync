package handlers

import (
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jakkayy/envSync/internal/database"
)

func PullSync(c *gin.Context) {
	projectID := c.Query("project_id")
	envName := c.Query("env")
	versionStr := c.Query("version")

	if projectID == "" || envName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id and env query params are required"})
		return
	}

	var env database.Environment
	err := database.DB.Where("project_id = ? AND name = ?", projectID, envName).First(&env).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "environment or project not found"})
		return
	}

	var versionRecord database.EnvVersion
	if versionStr != "" {
		vNum, err := strconv.Atoi(versionStr)
		if err == nil {
			database.DB.Where("environment_id = ? AND version = ?", env.ID, vNum).First(&versionRecord)
		}
	}

	if versionRecord.ID == 0 {
		err = database.DB.Where("environment_id = ?", env.ID).Order("version desc").First(&versionRecord).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "no configuration version found for environment"})
			return
		}
	}

	user, _ := c.Get("User")
	userStr, ok := user.(string)
	if !ok {
		userStr = "unknown-user"
	}

	database.DB.Create(&database.AuditLog{
		ProjectID: projectID,
		Action:    "PULL",
		EnvName:   envName,
		User:      userStr,
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	})

	encodedPayload := base64.StdEncoding.EncodeToString(versionRecord.EncryptedPayload)

	c.JSON(http.StatusOK, gin.H{
		"project_id": projectID,
		"env":        envName,
		"version":    versionRecord.Version,
		"payload":    encodedPayload,
		"message":    versionRecord.Message,
		"created_at": versionRecord.CreatedAt,
		"created_by": versionRecord.CreatedBy,
	})
}
