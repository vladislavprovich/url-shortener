package repository

import (
	"database/sql"
)

func InitDB(connStr string) (*sql.DB, error) {
	//connStr = "user=postgres password=password dbname=urlshortener host=db sslmode=disable host=localhost"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	// Verify connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
