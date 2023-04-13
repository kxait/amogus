package cracker

import (
	"amogus/child/state"
	"amogus/common"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func hashSha512(origin string) *common.HashPair {
	bytes := []byte(origin)
	sha := sha512.New()
	sha.Write(bytes)

	hash := sha.Sum(nil)
	result := &common.HashPair{
		Hash:   fmt.Sprintf("%x", hash[:]),
		Origin: origin,
	}

	return result
}

func hashSha256(origin string) *common.HashPair {
	bytes := []byte(origin)
	sha := sha256.New()
	sha.Write(bytes)

	hash := sha.Sum(nil)
	result := &common.HashPair{
		Hash:   fmt.Sprintf("%x", hash[:]),
		Origin: origin,
	}

	return result
}

func hashShadow(origin string, state *state.ChildState) *common.HashPair {
	e := (*state.ShadowCrypter).Crypt([]byte(origin)).String()

	return &common.HashPair{
		Hash:   e,
		Origin: origin,
	}
}
