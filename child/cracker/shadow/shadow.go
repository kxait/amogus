package shadow

import (
	"amogus/common"
	"regexp"
	"sync"

	"github.com/nathanaelle/password/v2"
)

const charset string = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func HashSha512Crypt(cr *password.Crypter, pass string) *common.HashPair {
	p := (*cr).Crypt([]byte(pass)).String()

	return &common.HashPair{
		Hash:   p,
		Origin: pass,
	}
}

func GetSaltySha512Crypter(line string) *password.Crypter {
	salt := ExtractSha512Salt(line)
	c, ok := password.SHA512.CrypterFound(salt)
	if !ok {
		panic("ok")
	}

	return &c
}

func ExtractSha512Salt(line string) string {
	saltRe, _ := regexp.Compile("\\$[a-z0-9]+\\$[a-zA-Z/.0-9]+\\$")
	return saltRe.FindString(line)
}

var once sync.Once
var cryptCharsetMap map[rune]byte

func getCharsetMap() *map[rune]byte {
	once.Do(func() {
		cryptCharsetMap = make(map[rune]byte)
		for i, c := range charset {
			cryptCharsetMap[c] = byte(i)
		}
	})

	return &cryptCharsetMap
}
