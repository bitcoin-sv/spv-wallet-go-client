package transports

import "errors"

// ErrAdminKey admin key not set
var ErrAdminKey = errors.New("an admin key must be set to be able to create an xpub")

// ErrNoClientSet is when no client is set
var ErrNoClientSet = errors.New("no transport client set")
