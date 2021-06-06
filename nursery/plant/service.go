package plant

// Service is the plant Service.
type Service struct {
	PlantRepository Repository
}

// CreatePlant creates a new Plant
// It returns Plant ID.
func (s Service) CreatePlant(p *Plant) (int64, error) {
	plantID, err := s.PlantRepository.CreatePlant(p)
	if err != nil {
		return 0, err
	}
	return plantID, nil
}

// GetPlant gets the Plant information from its ID
func (s Service) GetPlant(plantID string) (*Plant, error) {
	p, err := s.PlantRepository.FindPlantByID(plantID)
	if err != nil {
		return nil, err
	}
	return p, nil
}
