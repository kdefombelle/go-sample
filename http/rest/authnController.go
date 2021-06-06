package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/kdefombelle/go-sample/jwt"
	"github.com/kdefombelle/go-sample/logger"
	"github.com/kdefombelle/go-sample/nursery/account"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate = validator.New()

// LoginRequest the login request parameters.
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse the login request response.
type LoginResponse struct {
	UserID string `json:"userid" validate:"required"`
}

// AuthnController the Authentication Controller.
type AuthnController struct {
	AccountService *account.Service
	JwtService     jwt.Service
}

// Signin to authenticate user based on a LoginRequest.
// It returns a LoginResponse.
func (ac *AuthnController) Signin(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Info("Signin attempt")
	bytesUser, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Logger.Warnf("Cannot read login request [%v]", err)
		invalidRequest(w, err)
		return
	}
	request := &LoginRequest{}

	err = json.Unmarshal(bytesUser, request)
	logger.Logger.Infof("Signin attempt for [%s]", request.Username)
	if err != nil {
		logger.Logger.Warnw("Cannot deserialise", "request", request)
		invalidRequest(w, err)
		return
	}

	err = validate.Struct(request)
	if err != nil {
		logger.Logger.Warnw("Request not validated", "request", request)
		invalidRequest(w, err)
		return
	}

	// authenticating the user
	account, err := ac.AccountService.Authenticate(request.Username, request.Password)
	if err != nil {
		logger.Logger.Warnw("Cannot check user authorization", "err", err)
		invalidRequest(w, errors.New("invalid username or password"))
		return
	}

	//  creating the token
	token, err := ac.JwtService.NewToken(request.Username)
	if err != nil {
		logger.Logger.Warnw("Cannot create authorization", "err", err)
		invalidRequest(w, errors.New("cannot create authorization"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token.AccessToken,
		Expires:  token.ExpriresIn,
		HttpOnly: true,
		Secure:   false, //TODO: set to true once HTTPS (i.e. the cookie cannot be used out of HTTPS)
		Path:     "/",
		// SameSite: http.SameSiteStrictMode, //TODO: check if needed to add later when running on samesite - this is to avoid CSRF
	})

	//https://www.oauth.com/oauth2-servers/access-tokens/access-token-response/
	w.Header().Add("Cache-Control", "no-store")
	w.Header().Add("Pragma", "no-cache")

	response := LoginResponse{
		UserID: account.Username,
	}
	logger.Logger.Infof("Signin successfull for [%s]", request.Username)
	json, _ := json.Marshal(response)
	_, _ = w.Write(json)
}
