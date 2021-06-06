package authn

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

//Encrypter is the interface that wraps a basic Encrypty method.
type Encrypter interface {
	Encrypt(s string) string
}

//Md5Encrypter encrypts in Md5.
type Md5Encrypter struct {
}

//Encrypt a string.
func (e *Md5Encrypter) Encrypt(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	encryptedString := hex.EncodeToString(hasher.Sum(nil))
	fmt.Printf("%s", encryptedString)
	return encryptedString
}
