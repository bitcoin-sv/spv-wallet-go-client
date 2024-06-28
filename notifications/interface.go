package notifications

import "context"

type WebhookSubscriber interface {
	AdminSubscribeWebhook(ctx context.Context, webhookURL, tokenHeader, tokenValue string) error
	AdminUnsubscribeWebhook(ctx context.Context, webhookURL string) error
}
