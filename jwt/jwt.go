package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kdefombelle/go-sample/logger"
)

// Token is a JWT token.
type Token struct {
	Type        string
	AccessToken string
	ExpriresIn  time.Time
}

// Service exhibits the methods one can call on a jwt.Service.
// It has a PrivateKeyReader to read a private key from a location depending the implementation.
type Service struct {
	PrivateKeyReader PrivateKeyReader
}

// CustomClaims hold the extended claims for the application.
// i.e. standard claims and the username
type CustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// NewToken returns a new jwt.Token.
func (s Service) NewToken(username string) (*Token, error) {
	expiry := time.Now().UTC().Add(time.Minute * 15)
	// Set some claims
	claims := CustomClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiry.Unix(),
			Id:        "1",
			Issuer:    "nursery",
		},
	}
	// Create the token
	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)

	privateKeyBytes, err := s.PrivateKeyReader.Read()
	if err != nil {
		logger.Logger.Errorf("Cannot read private key [%s]", err)
		return &Token{}, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		logger.Logger.Errorf("Cannot parse private key [%s]", err)
		return &Token{}, errors.New("Cannot parse private key")
	}
	logger.Logger.Debugf("Private key parsed")

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(privateKey)
	logger.Logger.Debugf("Token signed")
	if err != nil {
		logger.Logger.Errorf("Cannot sign token", err)
		return &Token{}, errors.New("Cannot sign token")
	}
	jwtToken := Token{
		Type:        "Bearer",
		AccessToken: tokenString,
		ExpriresIn:  expiry,
	}
	logger.Logger.Debugf("Token created")
	return &jwtToken, nil

}

// ErrSignature is returned by jwt.Service.Verify method when the token signature is invalid.
var ErrSignature = errors.New("Token signature invalid")

// Verify a token.
// Returns false along with the associated error
//	if the private key cannot be read
//	if the key cannot be parsed
//	if the claims cannot be parsed
//	if the claims are not of type jwt.CustomClaims.
// Returns false and jwt.ErrSignature if the token signature is invalid.
// Returns true if the verification is successfull.
func (s Service) Verify(tokenString string) (bool, error) {
	tokenStringWithoutBearer := strings.TrimLeft(tokenString, "Bearer")
	tokenStringWithoutBearerAndSpace := strings.TrimSpace(tokenStringWithoutBearer)
	token, err := jwt.ParseWithClaims(tokenStringWithoutBearerAndSpace, &CustomClaims{}, func(*jwt.Token) (interface{}, error) {
		privateKeyBytes, err := s.PrivateKeyReader.Read()
		if err != nil {
			logger.Logger.Errorf("Cannot read private key [%s]", err)
			return false, err
		}

		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
		if err != nil {
			logger.Logger.Errorf("Cannot parse private key [%s]", err)
			return false, err
		}

		logger.Logger.Debugf("Private key parsed")
		return &privateKey.PublicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return false, ErrSignature
		}
		logger.Logger.Errorf("Could not parse claims [%s]", err)
		return false, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		logger.Logger.Error("Could not cast claims [%s]", ok)
		return false, nil
	}
	if claims.ExpiresAt < time.Now().UTC().Unix() {
		logger.Logger.Error("Authorization is expired")
		return false, nil
	}
	return true, nil
}

// String returns a string representation of the Token.
func (t Token) String() string {
	return fmt.Sprintf("%s %s", t.Type, t.AccessToken)
}
