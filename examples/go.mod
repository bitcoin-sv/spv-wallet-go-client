module github.com/bitcoin-sv/spv-wallet-go-client/examples

go 1.22.5

replace github.com/bitcoin-sv/spv-wallet-go-client => ../

require (
	github.com/bitcoin-sv/spv-wallet-go-client v0.0.0-00010101000000-000000000000
	github.com/bitcoin-sv/spv-wallet/models v1.0.0-beta.34
)

require (
	github.com/bitcoin-sv/go-sdk v1.1.14 // indirect
	github.com/boombuler/barcode v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pquerna/otp v1.4.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
)
