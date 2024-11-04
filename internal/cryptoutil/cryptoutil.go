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
	XpubKeyLength = 111
	ChainInternal = uint32(1)
	ChainExternal = uint32(0)
)

func Hash(s string) string {
	bb := sha256.Sum256([]byte(s))
	return hex.EncodeToString(bb[:])
}

func RandomHex(n int) (string, error) {
	bb := make([]byte, n)
	_, err := rand.Read(bb)
	if err != nil {
		return "", fmt.Errorf("failed to read bytes after rand: %w", err)
	}

	return hex.EncodeToString(bb), nil
}

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
	ErrMaxUint32LimitExceeded  = errors.New("max uint32 value exceeded")
	ErrNegativeValueNotAllowed = errors.New("negative value is not allowed")
	ErrHexHashPartIntParse     = errors.New("parse hex hash part to int64 failed")
)
