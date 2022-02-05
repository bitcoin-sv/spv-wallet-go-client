package transports

import "github.com/BitcoinSchema/xapi/bux"

// Recipients is a struct for recipients
type Recipients struct {
	To       string
	Satoshis uint64
	OpReturn *bux.OpReturn
}
