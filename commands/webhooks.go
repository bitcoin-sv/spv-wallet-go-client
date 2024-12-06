package commands

// CreateWebhookSubscription holds the arguments required to register a webhook subscription mechanism.
// This struct is used to define the details necessary for subscribing to webhook events.
type CreateWebhookSubscription struct {
	URL         string `json:"url"`         // The endpoint where webhook events will be sent. This must be a valid and reachable URL.
	TokenHeader string `json:"tokenHeader"` // The name of the HTTP header used for authentication in the subscription requests.
	TokenValue  string `json:"tokenValue"`  // The value of the authentication token that will be included in the TokenHeader.
}

// CancelWebhookSubscription holds the arguments required to cancel and remove a previously registered webhook subscription.
// This struct specifies the subscription endpoint that should be canceled.
type CancelWebhookSubscription struct {
	URL string `json:"url"` // The endpoint URL of the subscription to be removed. This must match the URL of an existing subscription.
}
