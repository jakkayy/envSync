package integration

import (
	"encoding/base64"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jakkayy/envSync/internal/database"
	"github.com/jakkayy/envSync/internal/server/handlers"
	"github.com/jakkayy/envSync/internal/server/middleware"
	"github.com/jakkayy/envSync/pkg/client"
)

func TestAPISyncFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("ENVSYNC_SKIP_AUTH", "true")

	db, err := database.InitDB("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	_ = database.AutoMigrate(db)

	router := gin.New()
	router.Use(middleware.LoggerAndRequestID())

	apiV1 := router.Group("/api/v1")
	apiV1.Use(middleware.AuthRequired())
	{
		apiV1.POST("/projects", handlers.CreateProject)
		apiV1.POST("/sync/push", handlers.PushSync)
		apiV1.GET("/sync/pull", handlers.PullSync)
		apiV1.GET("/projects/:id/history", handlers.GetProjectHistory)
	}

	ts := httptest.NewServer(router)
	defer ts.Close()

	apiClient := client.NewAPIClient(ts.URL, "")

	projID := "proj_test123"
	rawPayload := []byte("DB_HOST=localhost\nDB_PORT=5432\n")
	encPayload := base64.StdEncoding.EncodeToString(rawPayload)

	version, err := apiClient.Push(projID, "dev", encPayload, "Initial commit")
	if err != nil {
		t.Fatalf("Push failed: %v", err)
	}
	if version != 1 {
		t.Errorf("expected version 1, got %d", version)
	}

	pulledPayload, pulledVersion, err := apiClient.Pull(projID, "dev", "")
	if err != nil {
		t.Fatalf("Pull failed: %v", err)
	}
	if pulledVersion != 1 {
		t.Errorf("expected pulled version 1, got %d", pulledVersion)
	}
	if pulledPayload != encPayload {
		t.Errorf("pulled payload mismatch")
	}

	history, err := apiClient.GetHistory(projID, "dev")
	if err != nil {
		t.Fatalf("GetHistory failed: %v", err)
	}
	if len(history) != 1 {
		t.Errorf("expected 1 history item, got %d", len(history))
	}
}
