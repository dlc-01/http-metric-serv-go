package hashing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func HashingDate(key string, data []byte) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)

	return hex.EncodeToString(h.Sum(nil))
}

func CheckingHash(hash, key string, data []byte) (bool, error) {
	decodeHash, err := hex.DecodeString(hash)
	if err != nil {
		return false, fmt.Errorf("cannot decodeHash from string: %w", err)
	}

	newHash, err := hex.DecodeString(HashingDate(key, data))
	if err != nil {
		return false, fmt.Errorf("cannot encode from string: %w", err)
	}

	return hmac.Equal(newHash, decodeHash), nil
}
