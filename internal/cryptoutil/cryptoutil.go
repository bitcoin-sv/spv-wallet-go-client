package cryptoutil

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
)

const (
	// XpubKeyLength is the length of an xPub string key.
	XpubKeyLength = 111

	// ChainInternal internal chain num.
	ChainInternal = uint32(1)

	// ChainExternal external chain num.
	ChainExternal = uint32(0)
)

// Hash returns encoded string from SHA256 checksum of the s value.
func Hash(s string) string {
	bb := sha256.Sum256([]byte(s))
	return hex.EncodeToString(bb[:])
}

// RandomHex returns a random hexadecimal string of length `n`.
// An non-nil error is returned when random byte generation process failure occurs (rand.Read).
func RandomHex(n int) (string, error) {
	bb := make([]byte, n)
	_, err := rand.Read(bb)
	if err != nil {
		return "", fmt.Errorf("failed to read bytes after rand: %w", err)
	}
	return hex.EncodeToString(bb), nil
}

// DeriveChildKeyFromHex derives a child extended key from a BIP32 master key using a hexadecimal hash.
// The hexadecimal string is parsed into a 32-bit unsigned integers slice, each of which is used
// to derive a successive child key from the provided master or parent key.
func DeriveChildKeyFromHex(key *bip32.ExtendedKey, hexHash string) (*bip32.ExtendedKey, error) {
	nums, err := ParseChildNumsFromHex(hexHash)
	if err != nil {
		return nil, fmt.Errorf("failed to return parsed child nums from hex hash: %w", err)
	}

	child := key
	for _, n := range nums {
		child, err = child.Child(n)
		if err != nil {
			return nil, fmt.Errorf("failed to return derived child extended key: %w", err)
		}
	}
	return child, nil
}

// ParseChildNumsFromHex parses a hexadecimal string into multiple 32-bit unsigned integers.
// The input hex string is divided into 8-character chunks, each of which is interpreted as a 32-bit
// unsigned integer in hexadecimal format. The function returns a slice of these integers or an error
// if any part of the hex string is not valid.
func ParseChildNumsFromHex(hexHash string) ([]uint32, error) {
	if hexHash == "" {
		return nil, nil
	}

	const size = 8
	parts := (len(hexHash) + size - 1) / size // Avoids the need for floating-point division and ensures correct rounding up.
	var nums []uint32
	for i := 0; i < parts; i++ {
		start := i * size
		end := start + size
		if end > len(hexHash) {
			end = len(hexHash) // Adjust end to fit remaining substring.
		}
		num, err := parseHexPart(hexHash[start:end])
		if err != nil {
			return nil, fmt.Errorf("failed to parse hex part %q: %w", hexHash[start:end], err)
		}
		nums = append(nums, num)
	}
	return nums, nil
}

// parseHexPart converts a hexadecimal string to a uint32 value.
// The input string is expected to represent a valid hexadecimal number, and its value
// must fit within the range of a 32-bit unsigned integer.
func parseHexPart(part string) (uint32, error) {
	i, err := strconv.ParseInt(part, 16, 64)
	if err != nil {
		return 0, errors.Join(err, ErrHexHashPartIntParse)
	}
	u, err := Int64ToUint32(i % math.MaxInt32)
	if err != nil {
		return 0, fmt.Errorf("failed to convert int64 to uint32: %w", err)
	}
	return u, nil
}

// Int64ToUint32 converts an int64 value to a uint32, ensuring that the input is within the valid range for a uint32.
// The function performs a range check to ensure the int64 value is non-negative and does not exceed the maximum value
// for a uint32 (which is 2^32 - 1).
func Int64ToUint32(value int64) (uint32, error) {
	if value < 0 {
		return 0, ErrNegativeValueNotAllowed
	}
	if value > math.MaxUint32 {
		return 0, ErrMaxUint32LimitExceeded
	}
	return uint32(value), nil
}

var (
	// ErrMaxUint32LimitExceeded occurs when attempting to convert an int64 value that exceeds the maximum uint32 limit.
	ErrMaxUint32LimitExceeded = errors.New("max uint32 value exceeded")

	// ErrNegativeValueNotAllowed occurs when attempting to convert a negative int64 value to uint32.
	ErrNegativeValueNotAllowed = errors.New("negative value is not allowed")

	// ErrHexHashPartIntParse occurs when attempting to parse part of the hex hash to int64.
	ErrHexHashPartIntParse = errors.New("parse hex hash part to int64 failed")
)
