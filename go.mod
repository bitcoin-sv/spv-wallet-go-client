module github.com/BuxOrg/go-buxclient

go 1.16

replace github.com/BitcoinSchema/xapi => ../xapi

require (
	github.com/BitcoinSchema/xapi v0.0.0-00010101000000-000000000000
	github.com/bitcoinschema/go-bitcoin/v2 v2.0.0-alpha.2
	github.com/libsv/go-bk v0.1.6
	github.com/libsv/go-bt/v2 v2.1.0-beta.2.0.20211221142324-0d686850c5e0
	github.com/machinebox/graphql v0.2.2
	github.com/matryer/is v1.4.0 // indirect
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20220209154931-65fa2f7aa847 // indirect
)
