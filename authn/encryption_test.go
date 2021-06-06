package authn

import (
	"testing"

	"github.com/kdefombelle/go-sample/shared/test"
)

func TestMd5Encrypter_JohnDoe(t *testing.T) {
	encrypter := Md5Encrypter{}
	res := encrypter.Encrypt("johndoe")
	test.CheckString(t, "encrypted string", res, "6579e96f76baa00787a28653876c6127")
}

func TestMd5Encrypter_JaneDoe(t *testing.T) {
	encrypter := Md5Encrypter{}
	res := encrypter.Encrypt("janedoe")
	test.CheckString(t, "encrypted string", res, "a8c0d2a9d332574951a8e4a0af7d516f")
}
