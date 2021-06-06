package plant

import (
	"database/sql"
	"testing"

	"github.com/kdefombelle/go-sample/shared/test"
	_ "github.com/proullon/ramsql/driver"
)

func before(t *testing.T) (*DbRepository, func()) {
	batch := []string{
		//INT AUTO_INCREMENT PRIMARY KEY not supported for id
		//DECIMAL(17,2) not supported for price
		`CREATE TABLE plant (
			id INT PRIMARY KEY, 
			name varchar(255) NOT NULL,
			planted_date DATE,
			price NUMBER NOT NULL,
			reserved BOOLEAN
		  );`,
		`INSERT INTO plant (id, name, planted_date, price, reserved) VALUES(1, 'Orchid', '2021-01-01', '175.23', 0);`,
		`INSERT INTO plant (id, name, planted_date, price, reserved) VALUES(2, 'Banana tree', '2019-01-01', '98.99', 1);`,
	}
	db, err := sql.Open("ramsql", t.Name())
	if err != nil {
		t.Fatalf("sql.Open : Error : %s\n", err)
	}
	for _, b := range batch {
		_, err = db.Exec(b)
		if err != nil {
			t.Fatalf("sql.Exec: Error: %s\n", err)
		}
	}
	repository := DbRepository{
		Db: db,
	}
	close := func() {
		db.Close()
	}
	return &repository, close
}

func TestFindPlantByID(t *testing.T) {
	repository, close := before(t)
	defer close()

	p, err := repository.FindPlantByID("1")
	if err != nil {
		t.Fatalf("Unexpected error: [%s]", err)
	}
	test.CheckInt64(t, "id", 1, p.ID)
	test.CheckString(t, "name", "Orchid", p.Name)
	test.CheckTime(t, "planted_date", test.GetDate(t, "2021-01-01"), p.PlantedDate)
	test.CheckFloat(t, "price", 175.23, p.Price)
	test.CheckBool(t, "reserved", false, p.Reserved)
}

func TestCreatePlant(t *testing.T) {
	plantedDate := test.GetDate(t, "2021-06-01")
	plant := Plant{
		ID:          1, //normally not provided as id auto increment
		Name:        "Gomanis Gofacilus",
		PlantedDate: plantedDate,
		Price:       19.99,
		Reserved:    true,
	}
	repository, close := before(t)
	defer close()

	_, err := repository.CreatePlant(&plant)
	if err != nil {
		t.Fatalf("Unexpected error: [%s]", err)
	}
}
