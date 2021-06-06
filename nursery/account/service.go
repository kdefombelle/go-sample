package account

import (
	"fmt"

	"github.com/kdefombelle/go-sample/authn"
)

// Service is the User Service.
type Service struct {
	AccountRepository Repository
	Encrypter         authn.Encrypter
}

// Authenticate a user based on his username and passoword.
func (s Service) Authenticate(username string, password string) (*Account, error) {
	a, err := s.AccountRepository.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	encryptedPassword := s.Encrypter.Encrypt(password)
	if encryptedPassword != a.Password {
		return nil, fmt.Errorf("invalid password for user %q", username)
	}
	return a, nil
}
