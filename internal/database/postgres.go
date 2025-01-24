package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectAndMigrate(host, user, password, dbname, port string) error {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}

	DB = db

	query := `
    CREATE TABLE IF NOT EXISTS swift_codes (
      id SERIAL PRIMARY KEY,
      swift_code VARCHAR(11) NOT NULL UNIQUE,
      bank_name VARCHAR(255) NOT NULL,
      address VARCHAR(255) NOT NULL,
      country_iso2 VARCHAR(2) NOT NULL,
      country_name VARCHAR(100) NOT NULL,
      is_headquarter BOOLEAN NOT NULL
    );
    `
	_, err = DB.Exec(query)
	return err
}
