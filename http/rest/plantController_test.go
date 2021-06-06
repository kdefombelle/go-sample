package rest

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/kdefombelle/go-sample/nursery/plant"
	"github.com/kdefombelle/go-sample/shared/mock"
)

func TestAddPlant(t *testing.T) {
	uc := PlantController{
		PlantService: &plant.Service{
			PlantRepository: &mock.PlantRepository{},
		},
	}
	addPlantRequest := `{
							"id": "30",
							"name" :"Doe",
							"planted_date" :"2021-01-01",
							"Price" :"19.99",
							"Reserved" :"false"
						}`
	req, err := http.NewRequest("POST", "/plant", bytes.NewBuffer([]byte(addPlantRequest)))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(uc.AddPlant)
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got [%v] want [%v]",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"id":"1"}`
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got [%v] want [%v]",
			recorder.Body.String(), expected)
	}
}

func TestGetPlant(t *testing.T) {
	plantController := PlantController{
		PlantService: &plant.Service{
			PlantRepository: &mock.PlantRepository{},
		},
	}

	r := chi.NewRouter()
	r.Route("/plant", func(r chi.Router) {
		r.Get("/{id}", plantController.GetPlant)
	})

	testServer := httptest.NewServer(r)
	res, err := http.Get(testServer.URL + "/plant/30")
	if err != nil {
		t.Fatal(err)
	}

	if status := res.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got [%v] want [%v]",
			status, http.StatusOK)
	}

	expected := `{"id":"30","name":"MockPlant","planted_date":"2021-01-01","price":"99.99","reserved":"false"}`
	checkExpectedJSON(t, res, expected)
}

func checkExpectedJSON(t *testing.T, res *http.Response, expected string) {
	builder := new(strings.Builder)
	_, _ = io.Copy(builder, res.Body)
	if builder.String() != expected {
		t.Errorf("handler returned unexpected body: got [%v] want [%v]", builder.String(), expected)
	}
}
