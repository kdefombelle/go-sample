package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" //mysql driver leveraged by sql.DB
	"github.com/kdefombelle/go-sample/logger"
)

//ConnectionFactory is a factory to create DB connections.
type ConnectionFactory struct{}

//Parameters hold the necessary parameters to connect to a database.
type Parameters struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

//Create a new database handle DB pool of database connections.
func (f ConnectionFactory) Create(params Parameters) (*sql.DB, func(), error) { //"username:password@tcp(127.0.0.1:3306)/test?parseTime=true"
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", params.Username, params.Password, params.Host, params.Port, params.Name)
	logger.Logger.Infof("Connecting to db at %q", url)
	db, err := sql.Open("mysql", url)
	close := func() {
		db.Close()
	}
	return db, close, err
}
