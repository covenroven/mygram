package utils

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// DBValidateExists checks if the given value in a database already exists
func DBValidateExists(db *sqlx.DB, value interface{}, table string, col string) bool {
	query := fmt.Sprintf("SELECT id FROM %s WHERE %s = $1;", table, col)
	var id int
	row := db.QueryRow(query, value)
	if err := row.Scan(&id); err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return false
		}
	}

	return true
}
