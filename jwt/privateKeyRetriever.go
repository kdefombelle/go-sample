package jwt

import (
	"fmt"
	"io/ioutil"
)

// PrivateKeyReader depicts the interface to read a private key from a location.
type PrivateKeyReader interface {
	Read() ([]byte, error)
}

// FilePrivateKeyReader implements PrivateKeyReader.
type FilePrivateKeyReader struct {
	KeyPath string
}

// Read a private key from a local file
func (r FilePrivateKeyReader) Read() ([]byte, error) {
	fmt.Printf("Reading key from [%s]", r.KeyPath)
	return ioutil.ReadFile(r.KeyPath)
}
