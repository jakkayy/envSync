package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jakkayy/envSync/internal/database"
)

type HistoryResponseItem struct {
	ID        uint   `json:"id"`
	Version   int    `json:"version"`
	EnvName   string `json:"env_name"`
	CreatedBy string `json:"created_by"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

func GetProjectHistory(c *gin.Context) {
	projectID := c.Param("id")
	envName := c.Query("env")

	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project id is required"})
		return
	}

	var envs []database.Environment
	query := database.DB.Where("project_id = ?", projectID)
	if envName != "" {
		query = query.Where("name = ?", envName)
	}
	if err := query.Find(&envs).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project or environments not found"})
		return
	}

	envIDs := make([]uint, len(envs))
	envMap := make(map[uint]string)
	for i, e := range envs {
		envIDs[i] = e.ID
		envMap[e.ID] = e.Name
	}

	var versions []database.EnvVersion
	if len(envIDs) > 0 {
		database.DB.Where("environment_id IN ?", envIDs).Order("created_at desc").Limit(50).Find(&versions)
	}

	var history []HistoryResponseItem
	for _, v := range versions {
		history = append(history, HistoryResponseItem{
			ID:        v.ID,
			Version:   v.Version,
			EnvName:   envMap[v.EnvironmentID],
			CreatedBy: v.CreatedBy,
			Message:   v.Message,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	var auditLogs []database.AuditLog
	database.DB.Where("project_id = ?", projectID).Order("created_at desc").Limit(20).Find(&auditLogs)

	c.JSON(http.StatusOK, gin.H{
		"project_id": projectID,
		"history":    history,
		"audit_logs": auditLogs,
	})
}
