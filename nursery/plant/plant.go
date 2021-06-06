package plant

import "time"

// Plant model struct
type Plant struct {
	ID          int64
	Name        string
	PlantedDate time.Time
	Price       float64
	Reserved    bool
}
