package exampleutil

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/bitcoin-sv/spv-wallet-go-client/config"
)

// NewDefaultConfig returns a new instance of the default example configuration.
func NewDefaultConfig() config.Config {
	return config.New()
}

// PrettyPrint formats the provided JSON content with proper indentation
// to improve readability. It also displays a title, framed by two lines
// of `~` characters, for better visual presentation.
func PrettyPrint(title string, JSON any) {
	sep := strings.Repeat("~", 100)
	fmt.Println(sep)
	fmt.Println(title)
	fmt.Println(sep)

	res, err := json.MarshalIndent(JSON, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(res))
	fmt.Println()
}
