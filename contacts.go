package walletclient

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet-go-client/transports"
)

// AcceptContact accepts a contact by paymail
func (b *WalletClient) AcceptContact(ctx context.Context, paymail string) transports.ResponseError {
	return b.transport.AcceptContact(ctx, paymail)
}

// RejectContact rejects a contact by paymail
func (b *WalletClient) RejectContact(ctx context.Context, paymail string) transports.ResponseError {
	return b.transport.RejectContact(ctx, paymail)
}
