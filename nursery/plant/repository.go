package plant

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/kdefombelle/go-sample/logger"
)

// DbRepository is a Repository to operate a User encapsulating a Database connection.
type DbRepository struct {
	Db *sql.DB
}

// Repository interface for operating a Plant
type Repository interface {
	CreatePlant(p *Plant) (int64, error)
	FindPlantByID(plantID string) (*Plant, error)
}

// CreatePlant creates a Plant.
func (r *DbRepository) CreatePlant(p *Plant) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	price := strconv.FormatFloat(p.Price, 'f', 2, 64)
	res, err := r.Db.ExecContext(ctx, `insert into plant (
											name,
											planted_date,
											price,
											reserved
										)
										values (?,?,?,?)`, p.Name, p.PlantedDate, price, p.Reserved)

	if err != nil {
		return 0, err
	}
	logger.Logger.Infof("Plant [%s] inserted", p)
	return res.LastInsertId()
}

// FindPlantByID finds a Plant from its ID.
func (r *DbRepository) FindPlantByID(plantID string) (*Plant, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	query := `select
					id,
					name,
					planted_date,
					price,
					reserved
				from
					plant 
				where
					plant.id = ?`

	logQuery(query, map[string]string{
		"plantID": plantID,
	})
	row := r.Db.QueryRowContext(ctx, query, plantID)

	var plant Plant
	err := row.Scan(
		&plant.ID,
		&plant.Name,
		&plant.PlantedDate,
		&plant.Price,
		&plant.Reserved,
	)

	if err != nil {
		return &plant, err
	}
	return &plant, nil
}

func logQuery(query string, parameters map[string]string) {
	log.Printf("Query [%s]", minify(query))
	for k, v := range parameters {
		log.Printf("Parameters [%s]:[%s]", k, v)
	}
}

func minify(s string) string {
	s = strings.Replace(s, "\n", " ", -1)
	s = strings.Replace(s, "\t", "", -1)
	return s
}
