package jwt

import (
	"strings"
	"testing"
	"time"
)

func TestVerify(t *testing.T) {
	jwtService := Service{
		PrivateKeyReader: FilePrivateKeyReader{KeyPath: "../nursery.rsa"},
	}
	token, _ := jwtService.NewToken("johndoe")
	valid, err := jwtService.Verify(token.String())
	if err != nil {
		t.Fatalf("Error while verifying token [%s]", err)
	}

	if !valid {
		t.Error("The token is invalid")
	}
}

func TestNewJwtToken(t *testing.T) {
	jwtService := Service{
		PrivateKeyReader: FilePrivateKeyReader{KeyPath: "../nursery.rsa"},
	}
	var token *Token
	var err error

	if token, err = jwtService.NewToken("johndoe"); err != nil {
		t.Fatalf("Error while creating token [%s]", err)
	}

	if token.AccessToken == "" {
		t.Error("Token access token is empty")
	}

	if token.ExpriresIn.Unix() < time.Now().UTC().Unix() {
		t.Error("Token is expired")
	}

	if strings.ToLower(token.Type) != "bearer" {
		t.Errorf("Token type is incorrect [%s]", token.Type)
	}
}
