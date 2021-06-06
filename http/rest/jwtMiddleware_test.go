package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/kdefombelle/go-sample/jwt"
)

func startServer(t *testing.T) (*http.Request, func()) {
	r := chi.NewRouter()
	jwtMiddleware := JwtMiddleware{
		JwtService: jwt.Service{
			PrivateKeyReader: jwt.FilePrivateKeyReader{KeyPath: "../../nursery.rsa"},
		},
	}
	r.Use(jwtMiddleware.Check)
	r.Post("/test", func(w http.ResponseWriter, r *http.Request) {})

	ts := httptest.NewServer(r)

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/test", nil)

	if err != nil {
		t.Fatal(err)
	}
	close := func() {
		ts.Close()
	}
	return req, close
}

func TestCheckNoCookie(t *testing.T) {
	req, close := startServer(t)
	defer close()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("The reponse has status code [%d], expected was [%d]", resp.StatusCode, http.StatusUnauthorized)
	}
}

func TestCheckCookieEmpty(t *testing.T) {
	req, close := startServer(t)
	defer close()
	c := &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().UTC().Add(time.Minute * 15),
		HttpOnly: true,
		Secure:   false, //TODO: set to true once HTTPS (i.e. the cookie cannot be used out of HTTPS)

	}
	req.AddCookie(c)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("The reponse has status code [%d], expected was [%d]", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestCheckCookieOk(t *testing.T) {
	req, close := startServer(t)
	defer close()
	jwtService := jwt.Service{
		PrivateKeyReader: jwt.FilePrivateKeyReader{KeyPath: "../../nursery.rsa"},
	}
	token, _ := jwtService.NewToken("karim")
	c := &http.Cookie{
		Name: "token",
		//extra * at beginning
		Value:    token.AccessToken,
		Expires:  token.ExpriresIn,
		HttpOnly: true,
		Secure:   false,
	}
	req.AddCookie(c)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("The reponse has status code [%d], expected was [%d]", resp.StatusCode, http.StatusOK)
	}
}
