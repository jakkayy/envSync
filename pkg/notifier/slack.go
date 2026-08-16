package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SlackNotifier struct {
	WebhookURL string
}

func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{WebhookURL: webhookURL}
}

func (s *SlackNotifier) SendPushNotification(projectName, envName string, version int, user, message string) error {
	if s.WebhookURL == "" {
		return nil
	}

	payload := map[string]interface{}{
		"text": fmt.Sprintf("📢 *envSync Alert:* `%s` updated environment `%s` to version *v%d*", projectName, envName, version),
		"blocks": []map[string]interface{}{
			{
				"type": "section",
				"text": map[string]string{
					"type": "mrkdwn",
					"text": fmt.Sprintf("📢 *envSync Config Update Notification*\n*Project:* %s | *Env:* `%s` | *Version:* v%d", projectName, envName, version),
				},
			},
			{
				"type": "context",
				"elements": []map[string]string{
					{
						"type": "mrkdwn",
						"text": fmt.Sprintf("👤 *Updated By:* @%s | 🕒 *Time:* %s", user, time.Now().Format("2006-01-02 15:04:05")),
					},
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack payload: %w", err)
	}

	resp, err := http.Post(s.WebhookURL, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to send Slack notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("slack webhook returned status code %d", resp.StatusCode)
	}

	return nil
}
