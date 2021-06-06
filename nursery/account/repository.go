package account

import (
	"context"
	"database/sql"
	"time"

	"github.com/kdefombelle/go-sample/logger"
)

//A DbRepository is a Repository to operate an Account encapsulating a Database connection.
type DbRepository struct {
	Db *sql.DB
}

//A Repository interface for operating an Account.
type Repository interface {
	Create(a *Account) error
	FindByUsername(username string) (*Account, error)
}

//Create an Account.
func (r DbRepository) Create(a *Account) error {
	stmt, err := r.Db.Prepare(`insert into account (
									username,
									password,)
								values (?,?)`)
	if err != nil {
		return err
	}
	_, err2 := stmt.Exec(a.Username, a.Password)
	if err2 != nil {
		return err
	}
	logger.Logger.Debugf("Account for [%s] inserted", a.Username)
	return nil
}

//FindByUsername an Account.
func (r DbRepository) FindByUsername(username string) (*Account, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	row := r.Db.QueryRowContext(ctx, `select
											username,
											password
										from
											account
										where
											username=?`, username)
	var account Account
	err := row.Scan(&account.Username, &account.Password)
	if err != nil {
		return &account, err
	}
	return &account, nil
}
