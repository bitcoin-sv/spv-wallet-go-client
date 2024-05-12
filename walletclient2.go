// // Package walletclient is a Go client for interacting with Spv Wallet.
package walletclient

// import (
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
// 

// 	"github.com/bitcoinschema/go-bitcoin/v2"
// 	"github.com/libsv/go-bk/bec"
// 	"github.com/libsv/go-bk/bip32"
// 	"github.com/libsv/go-bk/wif"
// 	"github.com/pkg/errors"

// 	"github.com/bitcoin-sv/spv-wallet-go-client/transports"
// )

// // ClientOps are used for client options
// type ClientOps func(c *WalletClient)

// // // WalletClient is the spv wallet go client representation.
// // type WalletClient struct {
// // 	transports.TransportService
// // 	accessKey        *bec.PrivateKey
// // 	accessKeyString  string
// // 	transport        transports.TransportService
// // 	transportOptions []transports.ClientOps
// // 	xPriv            *bip32.ExtendedKey
// // 	xPrivString      string
// // 	xPub             *bip32.ExtendedKey
// // 	xPubString       string
// // }

// // New create a new wallet client
// func NewOld(opts ...ClientOps) (*WalletClient, error) {
// 	client := &WalletClient{}

// 	for _, opt := range opts {
// 		opt(client)
// 	}

// 	var err error
// 	if client.xPrivString != "" {
// 		if client.xPriv, err = bitcoin.GenerateHDKeyFromString(client.xPrivString); err != nil {
// 			return nil, err
// 		}
// 		if client.xPub, err = client.xPriv.Neuter(); err != nil {
// 			return nil, err
// 		}
// 	} else if client.xPubString != "" {
// 		client.xPriv = nil
// 		if client.xPub, err = bitcoin.GetHDKeyFromExtendedPublicKey(client.xPubString); err != nil {
// 			return nil, err
// 		}
// 	} else if client.accessKeyString != "" {
// 		client.xPriv = nil
// 		client.xPub = nil

// 		var privateKey *bec.PrivateKey
// 		var decodedWIF *wif.WIF
// 		if decodedWIF, err = wif.DecodeWIF(client.accessKeyString); err != nil {
// 			// try as a hex string
// 			var errHex error
// 			if privateKey, errHex = bitcoin.PrivateKeyFromString(client.accessKeyString); errHex != nil {
// 				return nil, errors.Wrap(err, errHex.Error())
// 			}
// 		} else {
// 			privateKey = decodedWIF.PrivKey
// 		}
// 		client.accessKey = privateKey
// 	} else {
// 		return nil, errors.New("no keys available")
// 	}

// 	transportOptions := make([]transports.ClientOps, 0)
// 	if client.xPriv != nil {
// 		transportOptions = append(transportOptions, transports.WithXPriv(client.xPriv))
// 		transportOptions = append(transportOptions, transports.WithXPub(client.xPub))
// 	} else if client.xPub != nil {
// 		transportOptions = append(transportOptions, transports.WithXPub(client.xPub))
// 	} else if client.accessKey != nil {
// 		transportOptions = append(transportOptions, transports.WithAccessKey(client.accessKey))
// 	}
// 	if len(client.transportOptions) > 0 {
// 		transportOptions = append(transportOptions, client.transportOptions...)
// 	}

// 	if client.transport, err = transports.NewTransport(transportOptions...); err != nil {
// 		return nil, err
// 	}

// 	client.TransportService = client.transport

// 	return client, nil
// }

// // SetAdminKey set the admin key to use to create new xpubs
// func (b *WalletClient) SetAdminKey(adminKeyString string) error {
// 	adminKey, err := bip32.NewKeyFromString(adminKeyString)
// 	if err != nil {
// 		return err
// 	}

// 	b.transport.SetAdminKey(adminKey)

// 	return nil
// }

// // SetSignRequest turn the signing of the http request on or off
// func (b *WalletClient) SetSignRequest(signRequest bool) {
// 	b.transport.SetSignRequest(signRequest)
// }

// // IsSignRequest return whether to sign all requests
// func (b *WalletClient) IsSignRequest() bool {
// 	return b.transport.IsSignRequest()
// }

// // GetTransport returns the current transport service
// func (b *WalletClient) GetTransport() *transports.TransportService {
// 	return &b.transport
// }
