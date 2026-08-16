package audit

import (
	"fmt"
	"time"

	"github.com/jakkayy/envSync/internal/database"
)

type AuditEvent struct {
	ProjectID string
	Action    string
	EnvName   string
	User      string
	IPAddress string
	UserAgent string
}

var eventQueue = make(chan AuditEvent, 500)

func init() {
	go startWorker()
}

func startWorker() {
	for event := range eventQueue {
		if database.DB == nil {
			continue
		}
		err := database.DB.Create(&database.AuditLog{
			ProjectID: event.ProjectID,
			Action:    event.Action,
			EnvName:   event.EnvName,
			User:      event.User,
			IPAddress: event.IPAddress,
			UserAgent: event.UserAgent,
			CreatedAt: time.Now(),
		}).Error

		if err != nil {
			fmt.Printf("AuditLog async write error: %v\n", err)
		}
	}
}

// LogAsync enqueues an audit event for non-blocking background persistence
func LogAsync(event AuditEvent) {
	select {
	case eventQueue <- event:
	default:
		fmt.Println("Warning: Audit log queue full, event dropped")
	}
}
