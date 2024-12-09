package exampleutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/config"
)

// NewDefaultConfig returns a new instance of the default example configuration.
func NewDefaultConfig() config.Config {
	return config.New()
}

// Print formats the title using the specified format and arguments, then prints the object.
func Print(format string, args ...any) {
	Printf(format, nil, "", 0, args...)
}

// Printf formats the title using the specified format and arguments, then marshals and prints the object with custom separators.
func Printf(format string, a any, separatorChar string, separatorLen int, args ...any) {
	// Default separator character and length if not provided
	if separatorChar == "" {
		separatorChar = "~"
	}
	if separatorLen <= 0 {
		separatorLen = 100
	}

	separator := strings.Repeat(separatorChar, separatorLen)
	var buf bytes.Buffer

	// Build the output
	buf.WriteString(separator + "\n")
	buf.WriteString(fmt.Sprintf(format, args...) + "\n")
	buf.WriteString(separator + "\n")

	// Marshal the object if provided
	if a != nil {
		res, err := json.MarshalIndent(a, "", "  ")
		if err != nil {
			log.Printf("Error marshaling data for '%s': %v", format, err)
			buf.WriteString("<error marshaling data>\n")
		} else {
			buf.WriteString(string(res) + "\n")
		}
	}

	// Print the buffer content
	fmt.Print(buf.String())
}

func RandomPaymail() string {
	seed := time.Now().UnixNano()
	n := rand.New(rand.NewSource(seed)).Intn(500)
	addr := fmt.Sprintf("john.doe.%dtest@4chain.test.com", n)
	return addr
}
