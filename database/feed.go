package database

import (
	"database/sql"
	"fmt"
)

func GetAllFeeds(db *sql.DB) ([]string, error) {
	getAllFeeds := "SELECT feed FROM %s"
	rows, err := db.Query(fmt.Sprintf(getAllFeeds, feedsTable))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feeds := []string{}
	for rows.Next() {
		var feed string
		if err := rows.Scan(&feed); err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}

func GetFeeds(db *sql.DB, feed string) ([]string, error) {
	getFeeds := "SELECT feed FROM %s WHERE LOWER(feed) = LOWER('%s')"
	rows, err := db.Query(fmt.Sprintf(getFeeds, feedsTable, feed))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feeds := []string{}
	for rows.Next() {
		var feed string
		if err := rows.Scan(&feed); err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feeds, nil
}

func InsertFeed(db *sql.DB, feed string) error {
	insertFeed := "INSERT INTO %s (feed) VALUES ('%s')"
	if _, err := db.Exec(fmt.Sprintf(insertFeed, feedsTable, feed)); err != nil {
		return err
	}
	return nil
}

func DeleteFeed(db *sql.DB, feed string) error {
	deleteFeed := "DELETE FROM %s WHERE LOWER(feed) = LOWER('%s')"
	if _, err := db.Exec(fmt.Sprintf(deleteFeed, feedsTable, feed)); err != nil {
		return err
	}
	return nil
}
