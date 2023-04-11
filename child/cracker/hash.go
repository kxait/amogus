package cracker

import (
	"crypto/sha512"
	"fmt"
)

func hashSha512(origin string) *HashPair {
	bytes := []byte(origin)
	sha := sha512.New()
	sha.Write(bytes)

	hash := sha.Sum(nil)
	result := &HashPair{
		Hash:   fmt.Sprintf("%x", hash[:]),
		Origin: origin,
	}

	return result
}
