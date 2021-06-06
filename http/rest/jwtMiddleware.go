package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/kdefombelle/go-sample/common"
	"github.com/kdefombelle/go-sample/jwt"
	"github.com/kdefombelle/go-sample/logger"
)

// JwtMiddleware encapsulates a service checking a JWT token.
type JwtMiddleware struct {
	JwtService jwt.Service
}

// Check a JWT token from a "token" Cookie.
func (m JwtMiddleware) Check(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				log.Printf("[%s]: no cookie, request non authorized", http.StatusText(http.StatusUnauthorized))
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			apiError := common.APIError{
				Error:            "invalid_client",
				ErrorDescription: fmt.Sprintf("%s", "Cannot retrieve authorization"),
			}
			json, _ := json.Marshal(apiError)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write(json)
			log.Printf("[%s]: cannot retrieve cookie token", http.StatusText(http.StatusUnauthorized))
			return
		}
		logger.Logger.Debugf("Checking for token [%s]", c.Value)
		if valid, err := m.JwtService.Verify("Bearer " + c.Value); err != nil || !valid {
			if errors.Is(err, jwt.ErrSignature) {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
