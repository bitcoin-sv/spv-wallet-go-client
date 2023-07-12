// Package utils contains utility functions for the wallet like hashes and crypto functions
package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math"
	"strconv"

	"github.com/libsv/go-bk/bip32"
)

const (
	// XpubKeyLength is the length of an xPub string key
	XpubKeyLength = 111

	// ChainInternal internal chain num
	ChainInternal = uint32(1)

	// ChainExternal external chain num
	ChainExternal = uint32(0)

	// MaxInt32 max integer for int32
	MaxInt32 = int64(1<<(32-1) - 1)
)

// Hash returns the sha256 hash of the data string
func Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// RandomHex returns a random hex string and error
func RandomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// DeriveChildKeyFromHex derive the child extended key from the hex string
func DeriveChildKeyFromHex(hdKey *bip32.ExtendedKey, hexHash string) (*bip32.ExtendedKey, error) {
	var childKey *bip32.ExtendedKey
	childKey = hdKey

	childNums, err := GetChildNumsFromHex(hexHash)
	if err != nil {
		return nil, err
	}

	for _, num := range childNums {
		if childKey, err = childKey.Child(num); err != nil {
			return nil, err
		}
	}

	return childKey, nil
}

// GetChildNumsFromHex get an array of uint32 numbers from the hex string
func GetChildNumsFromHex(hexHash string) ([]uint32, error) {
	strLen := len(hexHash)
	size := 8
	splitLength := int(math.Ceil(float64(strLen) / float64(size)))
	childNums := make([]uint32, 0)
	for i := 0; i < splitLength; i++ {
		start := i * size
		stop := start + size
		if stop > strLen {
			stop = strLen
		}
		num, err := strconv.ParseInt(hexHash[start:stop], 16, 64)
		if err != nil {
			return nil, err
		}
		if num > MaxInt32 {
			num = num - MaxInt32
		}
		childNums = append(childNums, uint32(num)) // todo: re-work to remove casting (possible cutoff)
	}

	return childNums, nil
}
