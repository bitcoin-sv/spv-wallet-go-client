package exampleutil

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/bitcoin-sv/spv-wallet-go-client/config"
)

var ExampleConfig = config.NewDefaultConfig("http://localhost:3003")

func Print(s string, a any) {
	fmt.Println(strings.Repeat("~", 100))
	fmt.Println(s)
	fmt.Println(strings.Repeat("~", 100))
	res, err := json.MarshalIndent(a, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(res))
}
