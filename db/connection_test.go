package db

import (
	"fmt"
	"strconv"
	"testing"
)

func TestOpen(t *testing.T) {
	params := Parameters{
		Name:     "nursery",
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "nursery",
		Password: "nursery",
	}
	connectionFactory := ConnectionFactory{}
	db, close, err := connectionFactory.Create(params)
	if err != nil {
		t.Fatal(err)
	}
	defer close()
	result, err := db.Exec("select * from plant")
	if err != nil {
		fmt.Printf("error while connecting to database [%s]", err)
	}
	lastResultID, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("last result id %s", strconv.FormatInt(lastResultID, 10))
	fmt.Printf("row(s) affected id %s", strconv.FormatInt(rowsAffected, 10))
}
