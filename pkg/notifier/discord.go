package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DiscordNotifier struct {
	WebhookURL string
}

func NewDiscordNotifier(webhookURL string) *DiscordNotifier {
	return &DiscordNotifier{WebhookURL: webhookURL}
}

func (d *DiscordNotifier) SendPushNotification(projectName, envName string, version int, user, message string) error {
	if d.WebhookURL == "" {
		return nil
	}

	payload := map[string]interface{}{
		"username": "envSync Bot",
		"embeds": []map[string]interface{}{
			{
				"title":       fmt.Sprintf("📢 Config Updated: %s [%s]", projectName, envName),
				"color":       3066993,
				"description": fmt.Sprintf("Environment configuration updated to version **v%d**", version),
				"fields": []map[string]interface{}{
					{"name": "Project", "value": projectName, "inline": true},
					{"name": "Environment", "value": envName, "inline": true},
					{"name": "Version", "value": fmt.Sprintf("v%d", version), "inline": true},
					{"name": "Updated By", "value": user, "inline": true},
					{"name": "Message", "value": message, "inline": false},
				},
				"timestamp": time.Now().Format(time.RFC3339),
			},
		},
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Discord payload: %w", err)
	}

	resp, err := http.Post(d.WebhookURL, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to send Discord notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("discord webhook returned status code %d", resp.StatusCode)
	}

	return nil
}
