package account

import (
	"database/sql"
	"testing"

	"github.com/kdefombelle/go-sample/shared/test"
	_ "github.com/proullon/ramsql/driver"
)

func before(t *testing.T) (*DbRepository, func()) {
	batch := []string{
		`CREATE TABLE account (
			username varchar(20) NOT NULL,
			password varchar(255) NOT NULL,
		  ) ;`,
		`INSERT INTO account (username, password) VALUES('johndoe', '6579e96f76baa00787a28653876c6127');`,
		`INSERT INTO account (username, password) VALUES('janedoe', 'a8c0d2a9d332574951a8e4a0af7d516f');`,
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
		repository.Db.Close()
	}
	return &repository, close
}

func TestCreate(t *testing.T) {
	repository, close := before(t)
	defer close()

	a := Account{
		Username: "username",
		Password: "password",
	}
	err := repository.Create(&a)
	if err != nil {
		t.Fatalf("Unexpected error: [%s]", err)
	}
}

func TestFindByUsername(t *testing.T) {
	repository, close := before(t)
	defer close()

	a, err := repository.FindByUsername("johndoe")
	if err != nil {
		t.Fatalf("Unexpected error: [%s]", err)
	}

	test.CheckString(t, "username", "johndoe", a.Username)
	test.CheckString(t, "password", "6579e96f76baa00787a28653876c6127", a.Password)
}
