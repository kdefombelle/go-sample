package mock

import (
	"time"

	"github.com/kdefombelle/go-sample/logger"
	"github.com/kdefombelle/go-sample/nursery/account"
	"github.com/kdefombelle/go-sample/nursery/plant"
)

// AccountRepository mock account.Repository.
type AccountRepository struct {
}

//Create an Account.
func (ur *AccountRepository) Create(a *account.Account) error {
	return nil
}

// FindByUsername finds an Account from a username.
func (ur *AccountRepository) FindByUsername(username string) (*account.Account, error) {
	return &account.Account{
		Username: "johndoe",
		Password: "6579e96f76baa00787a28653876c6127", //johndoe
	}, nil
}

// PlantRepository mock for plant.Repository.
type PlantRepository struct {
}

// CreatePlant create a plant.
func (pr *PlantRepository) CreatePlant(p *plant.Plant) (int64, error) {
	return 1, nil
}

// FindPlantByID finds a plant from its id.
func (pr *PlantRepository) FindPlantByID(id string) (*plant.Plant, error) {
	if id == "30" {
		p := createPlant()
		return p, nil
	}
	return nil, nil
}

func createPlant() *plant.Plant {
	plantedDate := getDate("2021-01-01")
	p := plant.Plant{
		ID:          30,
		Name:        "MockPlant",
		PlantedDate: plantedDate,
		Price:       99.99,
		Reserved:    false,
	}
	return &p
}

// Create a date in "2006-01-02" format
func getDate(date string) time.Time {
	const shortForm = "2006-01-02"
	time, err := time.Parse(shortForm, date)
	if err != nil {
		logger.Logger.Fatalf("Unexpected error: [%s]", err)
	}
	return time
}
