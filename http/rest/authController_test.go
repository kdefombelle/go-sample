package rest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kdefombelle/go-sample/authn"
	"github.com/kdefombelle/go-sample/jwt"
	"github.com/kdefombelle/go-sample/nursery/account"

	"github.com/kdefombelle/go-sample/shared/mock"
)

func TestSignin(t *testing.T) {
	ac := AuthnController{
		AccountService: &account.Service{
			AccountRepository: &mock.AccountRepository{},
			Encrypter:         &authn.Md5Encrypter{},
		},
		JwtService: jwt.Service{
			PrivateKeyReader: jwt.FilePrivateKeyReader{KeyPath: "../../nursery.rsa"},
		},
	}
	request := `{"username": "johndoe", "password":"johndoe"}`
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(request)))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.Signin)
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got [%v] want [%v]",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	response := `{"userid":"johndoe"}`
	if recorder.Body.String() != response {
		t.Errorf("handler returned unexpected body: got [%v] want [%v]",
			recorder.Body.String(), response)
	}
}
