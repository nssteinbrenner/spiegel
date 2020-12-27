package database

import (
	"database/sql"
	"fmt"

	"github.com/nssteinbrenner/spiegel/scrapers"
)

func GetAllHistory(db *sql.DB) ([]scrapers.DownloadInfo, error) {
	getAllHistory := "SELECT show, link, episode, quality, uploader FROM %s"
	rows, err := db.Query(fmt.Sprintf(getAllHistory, historyTable))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := []scrapers.DownloadInfo{}
	for rows.Next() {
		var show string
		var link string
		var episode string
		var quality int
		var uploader string

		if err := rows.Scan(&show, &link, &episode, &quality, &uploader); err != nil {
			return nil, err
		}
		downloadInfo := scrapers.DownloadInfo{
			Show:     show,
			Link:     link,
			Episode:  episode,
			Quality:  quality,
			Uploader: uploader,
		}
		history = append(history, downloadInfo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}

func GetHistory(db *sql.DB, d scrapers.DownloadInfo) ([]scrapers.DownloadInfo, error) {
	getHistory := fmt.Sprintf("SELECT show, link, episode, quality, uploader FROM %s WHERE '' = '' ", historyTable)
	if d.Show != "" {
		getHistory += fmt.Sprintf("AND LOWER(show) = LOWER('%s') ", d.Show)
	}
	if d.Link != "" {
		getHistory += fmt.Sprintf("AND LOWER(link) = LOWER('%s') ", d.Link)
	}
	if d.Episode != "" {
		getHistory += fmt.Sprintf("AND LOWER(episode) = LOWER('%s') ", d.Episode)
	}
	if d.Quality != 0 {
		getHistory += fmt.Sprintf("AND quality = %d ", d.Quality)
	}
	if d.Uploader != "" {
		getHistory += fmt.Sprintf("AND LOWER(uploader) = LOWER('%s') ", d.Uploader)
	}
	rows, err := db.Query(getHistory)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := []scrapers.DownloadInfo{}
	for rows.Next() {
		var show string
		var link string
		var episode string
		var quality int
		var uploader string

		if err := rows.Scan(&show, &link, &episode, &quality, &uploader); err != nil {
			return nil, err
		}
		downloadInfo := scrapers.DownloadInfo{
			Show:     show,
			Link:     link,
			Episode:  episode,
			Quality:  quality,
			Uploader: uploader,
		}
		history = append(history, downloadInfo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil

}

func InsertHistory(db *sql.DB, d scrapers.DownloadInfo) error {
	insertHistory := `
	INSERT INTO %s (
		show,
		link,
		episode,
		quality,
		uploader
	) VALUES (
		'%s',
		'%s',
		'%s',
		%d,
		'%s'
	)`
	query := fmt.Sprintf(
		insertHistory,
		historyTable,
		d.Show,
		d.Link,
		d.Episode,
		d.Quality,
		d.Uploader,
	)
	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}

func DeleteHistory(db *sql.DB, d scrapers.DownloadInfo) error {
	deleteHistory := fmt.Sprintf("DELETE FROM %s WHERE '' = '' ", historyTable)
	if d.Show != "" {
		deleteHistory += fmt.Sprintf("AND LOWER(show) = LOWER('%s') ", d.Show)
	}
	if d.Link != "" {
		deleteHistory += fmt.Sprintf("AND LOWER(link) = LOWER('%s') ", d.Link)
	}
	if d.Episode != "" {
		deleteHistory += fmt.Sprintf("AND LOWER(episode) = LOWER('%s') ", d.Episode)
	}
	if d.Quality != 0 {
		deleteHistory += fmt.Sprintf("AND quality = %d ", d.Quality)
	}
	if d.Uploader != "" {
		deleteHistory += fmt.Sprintf("AND LOWER(uploader) = LOWER('%s') ", d.Uploader)
	}
	if _, err := db.Exec(deleteHistory); err != nil {
		return err
	}
	return nil
}
