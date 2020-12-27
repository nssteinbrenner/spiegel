package database

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/nssteinbrenner/spiegel/scrapers"
)

func ProcessResults(db *sql.DB) ([]scrapers.DownloadInfo, error) {
	var results []scrapers.DownloadInfo

	showEpisodeMap, err := getShowEpisodeMap(db)
	if err != nil {
		return nil, err
	}

	for show, episodeSlice := range showEpisodeMap {
		uploaderCountMap, err := getUploaderCountMap(db, show)
		if err != nil {
			return nil, err
		}
		for _, episode := range episodeSlice {
			format := `
			WITH maxQuality(quality) AS (
				SELECT MAX(quality)
				FROM %s
				WHERE LOWER(show) = LOWER('%s')
				AND LOWER(episode) = LOWER('%s')
			) SELECT show, link, episode, quality, uploader, batch
			FROM %s
			WHERE LOWER(episode) = LOWER('%s')
			AND LOWER(show) = LOWER('%s')
			AND quality = (SELECT quality from maxQuality)
			AND NOT EXISTS (
				SELECT show, episode, quality
				FROM %s
				WHERE LOWER(show) = LOWER('%s')
				AND LOWER(episode) = LOWER('%s')
				AND quality >= (SELECT quality from maxQuality)
			)`
			query := fmt.Sprintf(
				format,
				resultsTable,
				show,
				episode,
				resultsTable,
				episode,
				show,
				historyTable,
				show,
				episode,
			)

			rows, err := db.Query(query)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			episodeResults := []scrapers.DownloadInfo{}
			for rows.Next() {
				var downloadInfo scrapers.DownloadInfo

				err = rows.Scan(
					&downloadInfo.Show,
					&downloadInfo.Link,
					&downloadInfo.Episode,
					&downloadInfo.Quality,
					&downloadInfo.Uploader,
					&downloadInfo.Batch,
				)
				if err != nil {
					return nil, err
				}

				episodeResults = append(episodeResults, downloadInfo)
			}
			if err := rows.Err(); err != nil {
				return nil, err
			}

			if len(episodeResults) > 0 {
				results = append(results, filterResults(episodeResults, uploaderCountMap))
			}
		}
	}

	return results, nil
}

func filterResults(results []scrapers.DownloadInfo, uploaderCountMap map[string]int) scrapers.DownloadInfo {
	if len(results) >= 2 {
		results = sortResults(results, uploaderCountMap)
	}

	for _, result := range results {
		if result.Batch {
			return result
		}
	}

	return results[0]
}

func sortResults(results []scrapers.DownloadInfo, uploaderCountMap map[string]int) []scrapers.DownloadInfo {
	if len(results) < 2 {
		return results
	}
	left, right := 0, len(results)-1
	pivot := rand.Int() % len(results)

	results[pivot], results[right] = results[right], results[pivot]
	for i, result := range results {
		if uploaderCountMap[result.Uploader] > uploaderCountMap[results[right].Uploader] {
			results[left], results[i] = results[i], results[left]
			left++
		}
	}

	results[left], results[right] = results[right], results[left]

	sortResults(results[:left], uploaderCountMap)
	sortResults(results[left+1:], uploaderCountMap)

	return results
}

func getShowEpisodeMap(db *sql.DB) (map[string][]string, error) {
	showEpisodeMap := make(map[string][]string)
	results, err := GetAllResults(db)
	if err != nil {
		return nil, err
	}

	for _, downloadInfo := range results {
		showEpisodeMap[downloadInfo.Show] = append(showEpisodeMap[downloadInfo.Show], downloadInfo.Episode)
	}

	for show, episodeSlice := range showEpisodeMap {
		encountered := make(map[string]bool)
		dedupedEpisodeSlice := []string{}
		for _, ep := range episodeSlice {
			if _, ok := encountered[ep]; !ok {
				dedupedEpisodeSlice = append(dedupedEpisodeSlice, ep)
				encountered[ep] = true
			}
		}
		showEpisodeMap[show] = dedupedEpisodeSlice
	}

	return showEpisodeMap, nil
}

func getUploaderCountMap(db *sql.DB, show string) (map[string]int, error) {
	uploaderCountMap := make(map[string]int)
	uploaderQuery := `
	SELECT
		uploader,
		count(uploader) AS uploader_count
	FROM %s
	WHERE LOWER(show) = LOWER('%s')
	GROUP BY uploader`
	rows, err := db.Query(fmt.Sprintf(uploaderQuery, resultsTable, show))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var uploader string
		var count int

		err = rows.Scan(&uploader, &count)
		if err != nil {
			return nil, err
		}

		uploaderCountMap[uploader] = count
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return uploaderCountMap, nil
}
