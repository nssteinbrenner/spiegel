package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	databaseName = "spiegel"
	historyTable = "history"
	resultsTable = "results"
	feedsTable   = "feeds"
	showsTable   = "shows"
	qualityTable = "quality"
)

var defaultFeeds = []string{
	"nyaa.si",
	"horriblesubs.info",
	"tokyotosho.info",
	"animetosho.org",
}

var defaultQuality = []int{
	1080,
	720,
	480,
}

func Open(dbPath string) (*sql.DB, error) {
	databaseFile := fmt.Sprintf("%s.db", databaseName)
	if dbPath != "" {
		databaseFile = fmt.Sprintf("%s/%s.db", strings.TrimRight(dbPath, "/"), databaseName)
	}
	_, err := os.Stat(databaseFile)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(databaseFile)
			if err != nil {
				return nil, err
			}
			file.Close()
		}
	}

	spiegelDB, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared", databaseFile))
	if err != nil {
		return nil, err
	}
	spiegelDB.SetMaxOpenConns(1)
	if err := initDB(spiegelDB); err != nil {
		return nil, err
	}

	return spiegelDB, nil
}

func initDB(db *sql.DB) error {
	createDownloadTable := `
	CREATE %s TABLE IF NOT EXISTS %s(
		show TEXT,
		link TEXT,
		episode TEXT,
		quality INT,
		uploader TEXT
	)`
	createResultsTable := `
	CREATE %s TABLE IF NOT EXISTS %s(
		show TEXT,
		link TEXT,
		episode TEXT,
		quality INT,
		uploader TEXT,
		batch BOOLEAN
	)`
	createFeedsTable := `
	CREATE TABLE IF NOT EXISTS %s(
		feed TEXT PRIMARY KEY
	)`
	createShowsTable := `
	CREATE TABLE IF NOT EXISTS %s(
		show TEXT PRIMARY KEY
	)`
	createQualityTable := `
	CREATE TABLE IF NOT EXISTS %s(
		quality INT PRIMARY KEY
	)`
	if _, err := db.Exec(fmt.Sprintf(createDownloadTable, "", historyTable)); err != nil {
		return err
	}
	if _, err := db.Exec(fmt.Sprintf(createResultsTable, "TEMPORARY", resultsTable)); err != nil {
		return err
	}
	if _, err := db.Exec(fmt.Sprintf(createFeedsTable, feedsTable)); err != nil {
		return err
	}
	if _, err := db.Exec(fmt.Sprintf(createShowsTable, showsTable)); err != nil {
		return err
	}
	if _, err := db.Exec(fmt.Sprintf(createQualityTable, qualityTable)); err != nil {
		return err
	}

	configuredFeeds, err := GetAllFeeds(db)
	if err != nil {
		return err
	}
	if len(configuredFeeds) == 0 {
		for _, feed := range defaultFeeds {
			if err := InsertFeed(db, feed); err != nil {
				return err
			}
		}
	}
	configuredQualities, err := GetAllQualities(db)
	if err != nil {
		return err
	}
	if len(configuredQualities) == 0 {
		for _, quality := range defaultQuality {
			if err := InsertQuality(db, quality); err != nil {
				return err
			}
		}
	}

	return nil
}
