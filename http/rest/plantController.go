package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/kdefombelle/go-sample/logger"
	"github.com/kdefombelle/go-sample/nursery/plant"
)

// Plant contains the base info for plant to be created
type Plant struct {
	ID          string `json:"id"`
	Name        string `json:"name" validate:"required,max=50"`
	PlantedDate string `json:"planted_date" validate:"required,datetime=2006-01-02"`
	Price       string `json:"price" validate:"required,max=20"`
	Reserved    string `json:"reserved" validate:"required"`
}

// AddPlantRequest the parameter to be passed to add a Plant.
type AddPlantRequest struct {
	Plant
}

// AddPlantResponse the response to an addition of a Plant.
type AddPlantResponse struct {
	ID string `json:"id"`
}

// PlantController HTTP requests controller.
type PlantController struct {
	PlantService *plant.Service
}

const shortForm = "2006-01-02"

// GetPlant HTTP endpoint.
func (pc *PlantController) GetPlant(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Info("Getting plant information")
	ID := chi.URLParam(r, "id")

	p, err := pc.PlantService.GetPlant(ID)
	if err != nil {
		logger.Logger.Fatalf("Cannot get plant [%v]", err)
		internalServer(w, fmt.Errorf("cannot fetch plant id %q", ID))
		return
	}
	if p == nil {
		notFound(w, fmt.Errorf("%q not found", ID))
		return
	}

	plant := Plant{
		ID:          strconv.FormatInt(p.ID, 10),
		Name:        p.Name,
		PlantedDate: p.PlantedDate.Format(shortForm),
		Price:       strconv.FormatFloat(p.Price, 'f', 2, 64),
		Reserved:    strconv.FormatBool(p.Reserved),
	}

	response := plant

	json, _ := json.Marshal(response)
	_, _ = w.Write(json)
}

// AddPlant HTTP endpoint.
func (pc *PlantController) AddPlant(w http.ResponseWriter, r *http.Request) {
	logger.Logger.Infof("Creating new plant")

	request := &AddPlantRequest{}
	err := validateRequest(r.Body, request)
	if err != nil {
		logger.Logger.Errorf("Invalid add plant request [%v]", err)
		invalidRequest(w, errors.New("invalid request"))
		return
	}

	// insert
	plantedDate, _ := time.Parse(shortForm, request.PlantedDate)
	price, _ := strconv.ParseFloat(request.Price, 64)
	reserved, _ := strconv.ParseBool(request.Reserved)
	p := plant.Plant{
		Name:        request.Name,
		PlantedDate: plantedDate,
		Price:       price,
		Reserved:    reserved,
	}
	ID, err := pc.PlantService.CreatePlant(&p)
	if err != nil {
		logger.Logger.Fatalf("Cannot create plant [%v]", err)
		internalServer(w, errors.New("cannot create plant"))
		return
	}
	logger.Logger.Infof("New plant [%s] created", ID)

	response := AddPlantResponse{ID: strconv.FormatInt(ID, 10)}
	json, _ := json.Marshal(response)
	_, _ = w.Write(json)
}
