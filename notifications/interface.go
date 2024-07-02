package notifications

import "context"

// WebhookSubscriber - interface for subscribing and unsubscribing to webhooks
type WebhookSubscriber interface {
	AdminSubscribeWebhook(ctx context.Context, webhookURL, tokenHeader, tokenValue string) error
	AdminUnsubscribeWebhook(ctx context.Context, webhookURL string) error
}
