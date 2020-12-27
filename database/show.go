package database

import (
	"database/sql"
	"fmt"
)

func GetAllShows(db *sql.DB) ([]string, error) {
	getAllShows := "SELECT show FROM %s"
	rows, err := db.Query(fmt.Sprintf(getAllShows, showsTable))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	shows := []string{}
	for rows.Next() {
		var show string
		if err := rows.Scan(&show); err != nil {
			return nil, err
		}
		shows = append(shows, show)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return shows, nil
}

func GetShows(db *sql.DB, show string) ([]string, error) {
	getShows := "SELECT show FROM %s WHERE LOWER(show) = LOWER('%s')"
	rows, err := db.Query(fmt.Sprintf(getShows, showsTable, show))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	shows := []string{}
	for rows.Next() {
		var show string
		if err := rows.Scan(&show); err != nil {
			return nil, err
		}
		shows = append(shows, show)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return shows, nil
}

func InsertShow(db *sql.DB, show string) error {
	insertShow := "INSERT INTO %s (show) VALUES ('%s')"
	if _, err := db.Exec(fmt.Sprintf(insertShow, showsTable, show)); err != nil {
		return err
	}
	return nil
}

func DeleteShow(db *sql.DB, show string) error {
	deleteShow := "DELETE FROM %s WHERE LOWER(show) = LOWER('%s')"
	if _, err := db.Exec(fmt.Sprintf(deleteShow, showsTable, show)); err != nil {
		return err
	}
	return nil
}
