package run

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nssteinbrenner/spiegel/config"
	"github.com/nssteinbrenner/spiegel/database"
	"github.com/nssteinbrenner/spiegel/metrics"
	"github.com/nssteinbrenner/spiegel/scrapers"
	"github.com/nssteinbrenner/spiegel/torrent"
	"github.com/sirupsen/logrus"
)

func StartRun(runConfig config.Config, feeds, shows []string, quality []int, logger *logrus.Entry) error {
	success := 1
	showsDownloaded := 0
	startTime := time.Now()
	defer func() {
		endTime := time.Now()
		duration := endTime.Sub(startTime).Seconds()
		if int(duration) < 0 {
			duration = 0
		}
		metrics.RunDuration.WithLabelValues(
			strconv.Itoa(showsDownloaded),
			strconv.Itoa(success),
		).Observe(duration)
	}()

	dbConn, err := database.Open(runConfig.DatabaseDirectory)
	if err != nil {
		logger.WithError(err).Info("Failed to open database")
		success = 0
		return err
	}
	defer dbConn.Close()
	if len(shows) == 0 {
		shows, err = database.GetAllShows(dbConn)
		if err != nil {
			logger.WithError(err).Info("Failed to get shows")
			success = 0
			return err
		}
	}
	if len(feeds) == 0 {
		feeds, err = database.GetAllFeeds(dbConn)
		if err != nil {
			logger.WithError(err).Info("Failed to get feeds")
			success = 0
			return err
		}
	}
	if len(quality) == 0 {
		quality, err = database.GetAllQualities(dbConn)
		if err != nil {
			logger.WithError(err).Info("Failed to get quality")
			success = 0
			return err
		}
	}

	scrapedTorrents := scrape(shows, feeds, quality, logger)
	for _, v := range scrapedTorrents {
		if err := database.InsertResults(dbConn, v); err != nil {
			logger.WithError(err).Info("Failed to insert results")
		}
	}
	res, err := database.ProcessResults(dbConn)
	if err != nil {
		logger.WithError(err).Info("Failed to select results")
		success = 0
		return err
	}
	tbt, err := torrent.GetTransmissionConnection(
		runConfig.TransmissionHost,
		runConfig.TransmissionPort,
		runConfig.TransmissionUser,
		runConfig.TransmissionPassword,
		runConfig.TransmissionHTTPS,
	)
	if err != nil {
		logger.WithError(err).Info("Failed to connect to transmission")
		success = 0
		return err
	}
	for _, v := range res {
		withoutLink := v
		withoutLink.Link = ""
		hist, err := database.GetHistory(dbConn, withoutLink)
		if err != nil {
			logger.WithError(err).Info("Failed to get history")
		}
		if len(hist) > 0 {
			continue
		}
		downloadDirectory := fmt.Sprintf("%s/%s", strings.TrimRight(runConfig.DownloadDirectory, "/"), v.Show)
		err = torrent.AddTransmissionTorrent(tbt, downloadDirectory, v.Link)
		if err != nil {
			logger.WithError(err).Info("Failed to add torrent")
			continue
		}
		err = database.InsertHistory(dbConn, v)
		if err != nil {
			logger.WithError(err).Info("Failed to insert history")
			continue
		}
		metrics.TotalEpisodesByShow.WithLabelValues(v.Show).Inc()
		showsDownloaded++
	}

	return nil
}

func scrape(shows []string, feeds []string, quality []int, logger *logrus.Entry) []scrapers.DownloadInfo {
	var scrapedTorrents []scrapers.DownloadInfo
	ret := make(chan scrapers.DownloadInfo)

	go func() {
		defer close(ret)
		scrapeFeeds(shows, feeds, quality, ret, logger)
	}()

	for i := range ret {
		scrapedTorrents = append(scrapedTorrents, i)
	}

	return scrapedTorrents
}

func scrapeFeeds(shows []string, feeds []string, quality []int, ret chan<- scrapers.DownloadInfo, logger *logrus.Entry) {
	logErr := func(feed, show string, err error, logger *logrus.Entry) {
		succeeded := 1
		if err != nil {
			succeeded = 0
			logger.WithError(err).WithFields(logrus.Fields{
				"feed": feed,
				"show": show,
			}).Error("Error encountered scraping feeds")
		}
		metrics.TotalScrapes.WithLabelValues(feed, strconv.Itoa(succeeded)).Inc()
	}

	var wgScrape sync.WaitGroup
	var wgRead sync.WaitGroup
	for _, feed := range feeds {
		for _, show := range shows {
			wgScrape.Add(1)
			go func(ret chan<- scrapers.DownloadInfo, feed string, show string, quality []int) {
				defer wgScrape.Done()
				c := make(chan scrapers.DownloadInfo)
				wgRead.Add(1)
				go func() {
					defer wgRead.Done()
					for i := range c {
						ret <- i
					}
				}()
				switch strings.ToLower(feed) {
				case "nyaa.si":
					go logErr(feed, show, scrapers.ScrapeNyaa(c, show, quality), logger)
				case "animetosho.org":
					go logErr(feed, show, scrapers.ScrapeAnimeTosho(c, show, quality), logger)
				case "tokyotosho.info":
					go logErr(feed, show, scrapers.ScrapeTokyoTosho(c, show, quality), logger)
				default:
					close(c)
					return
				}
			}(ret, feed, show, quality)
			// some sites (nyaa) will return a 429 if the requests come in too quickly.
			wgScrape.Wait()
		}
	}
	wgRead.Wait()
}
