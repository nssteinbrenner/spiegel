package database

import (
	"database/sql"
	"fmt"

	"github.com/nssteinbrenner/spiegel/scrapers"
)

func GetAllResults(db *sql.DB) ([]scrapers.DownloadInfo, error) {
	getAllResults := "SELECT show, link, episode, quality, uploader FROM %s"
	rows, err := db.Query(fmt.Sprintf(getAllResults, resultsTable))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []scrapers.DownloadInfo{}
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
		results = append(results, downloadInfo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func GetResults(db *sql.DB, d scrapers.DownloadInfo) ([]scrapers.DownloadInfo, error) {
	getResults := fmt.Sprintf("SELECT show, link, episode, quality, uploader FROM %s WHERE '' = ''", resultsTable)
	if d.Show != "" {
		getResults += fmt.Sprintf(" AND LOWER(show) = LOWER('%s')", d.Show)
	}
	if d.Link != "" {
		getResults += fmt.Sprintf(" AND LOWER(link) = LOWER('%s')", d.Link)
	}
	if d.Episode != "" {
		getResults += fmt.Sprintf(" AND LOWER(episode) = LOWER('%s')", d.Episode)
	}
	if d.Quality != 0 {
		getResults += fmt.Sprintf(" AND quality = %d", d.Quality)
	}
	if d.Uploader != "" {
		getResults += fmt.Sprintf(" AND LOWER(uploader) = LOWER('%s')", d.Uploader)
	}
	rows, err := db.Query(getResults)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []scrapers.DownloadInfo{}
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
		results = append(results, downloadInfo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil

}

func InsertResults(db *sql.DB, d scrapers.DownloadInfo) error {
	insertResults := `
	INSERT INTO %s (
		show,
		link,
		episode,
		quality,
		uploader,
		batch
	) VALUES (
		'%s',
		'%s',
		'%s',
		%d,
		'%s',
		%t
	)`
	query := fmt.Sprintf(
		insertResults,
		resultsTable,
		d.Show,
		d.Link,
		d.Episode,
		d.Quality,
		d.Uploader,
		d.Batch,
	)
	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}

func DeleteResults(db *sql.DB, d scrapers.DownloadInfo) error {
	deleteResults := fmt.Sprintf("DELETE FROM %s WHERE '' = ''", resultsTable)
	if d.Show != "" {
		deleteResults += fmt.Sprintf(" AND LOWER(show) = LOWER('%s')", d.Show)
	}
	if d.Link != "" {
		deleteResults += fmt.Sprintf(" AND LOWER(link) = LOWER('%s')", d.Link)
	}
	if d.Episode != "" {
		deleteResults += fmt.Sprintf(" AND LOWER(episode) = LOWER('%s')", d.Episode)
	}
	if d.Quality != 0 {
		deleteResults += fmt.Sprintf(" AND quality = %d", d.Quality)
	}
	if d.Uploader != "" {
		deleteResults += fmt.Sprintf(" AND LOWER(uploader) = LOWER('%s')", d.Uploader)
	}
	if _, err := db.Exec(deleteResults); err != nil {
		return err
	}
	return nil
}
