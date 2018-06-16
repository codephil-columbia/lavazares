package models

import (
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// var userSession *

//InitDB initialized connected to database
func InitDB(dataSourceName string) error {
	var err error
	db, err = sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return err
	}

	return nil
}
