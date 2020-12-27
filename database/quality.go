package database

import (
	"database/sql"
	"fmt"
)

func GetAllQualities(db *sql.DB) ([]int, error) {
	getAllQualities := "SELECT quality FROM %s"
	rows, err := db.Query(fmt.Sprintf(getAllQualities, qualityTable))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	qualities := []int{}
	for rows.Next() {
		var quality int
		if err := rows.Scan(&quality); err != nil {
			return nil, err
		}
		qualities = append(qualities, quality)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return qualities, nil
}

func GetQualities(db *sql.DB, quality int) ([]int, error) {
	getQualities := "SELECT quality FROM %s WHERE quality = %d"
	rows, err := db.Query(fmt.Sprintf(getQualities, qualityTable, quality))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	qualities := []int{}
	for rows.Next() {
		var quality int
		if err := rows.Scan(&quality); err != nil {
			return nil, err
		}
		qualities = append(qualities, quality)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return qualities, nil
}

func InsertQuality(db *sql.DB, quality int) error {
	insertFeed := "INSERT INTO %s (quality) VALUES (%d)"
	if _, err := db.Exec(fmt.Sprintf(insertFeed, qualityTable, quality)); err != nil {
		return err
	}
	return nil
}

func DeleteQuality(db *sql.DB, quality int) error {
	deleteFeed := "DELETE FROM %s WHERE quality = %d"
	if _, err := db.Exec(fmt.Sprintf(deleteFeed, qualityTable, quality)); err != nil {
		return err
	}
	return nil
}
