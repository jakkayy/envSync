package handlers

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jakkayy/envSync/internal/database"
)

type PushSyncRequest struct {
	ProjectID string `json:"project_id" binding:"required"`
	EnvName   string `json:"env_name" binding:"required"`
	Payload   string `json:"payload" binding:"required"`
	Message   string `json:"message"`
}

func PushSync(c *gin.Context) {
	var req PushSyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payloadBytes, err := base64.StdEncoding.DecodeString(req.Payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid base64 payload"})
		return
	}

	var env database.Environment
	err = database.DB.Where("project_id = ? AND name = ?", req.ProjectID, req.EnvName).First(&env).Error
	if err != nil {
		env = database.Environment{
			ProjectID:     req.ProjectID,
			Name:          req.EnvName,
			LatestVersion: 0,
		}
		database.DB.Create(&env)
	}

	newVersionNumber := env.LatestVersion + 1
	env.LatestVersion = newVersionNumber
	database.DB.Save(&env)

	user, _ := c.Get("User")
	userStr, ok := user.(string)
	if !ok {
		userStr = "unknown-user"
	}

	versionRecord := database.EnvVersion{
		EnvironmentID:    env.ID,
		Version:          newVersionNumber,
		EncryptedPayload: payloadBytes,
		CreatedBy:        userStr,
		Message:          req.Message,
	}
	database.DB.Create(&versionRecord)

	database.DB.Create(&database.AuditLog{
		ProjectID: req.ProjectID,
		Action:    "PUSH",
		EnvName:   req.EnvName,
		User:      userStr,
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "environment updated successfully",
		"version": newVersionNumber,
		"env":     req.EnvName,
	})
}
