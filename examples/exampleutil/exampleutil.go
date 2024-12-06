package exampleutil

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"math/rand"

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

func RandomPaymail() string {
	seed := time.Now().UnixNano()
	n := rand.New(rand.NewSource(seed)).Intn(500)
	addr := fmt.Sprintf("john.doe.%dtest@4chain.test.com", n)
	return addr
}
