package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nssteinbrenner/spiegel/config"
	"github.com/nssteinbrenner/spiegel/database"
	"github.com/nssteinbrenner/spiegel/run"
	"github.com/nssteinbrenner/spiegel/scrapers"
	"github.com/nssteinbrenner/spiegel/utils"
)

func handleFeeds(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	logger := utils.InitLogger()

	q := r.URL.Query()
	switch r.Method {
	case "POST":
		for _, i := range q["feed"] {
			err := database.InsertFeed(db, i)
			if err != nil {
				logger.WithError(err).Info("Failed to insert feed")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error() + "\n"))
				return
			}
		}
		w.Write([]byte(http.StatusText(http.StatusOK) + "\n"))
	case "DELETE":
		for _, i := range q["feed"] {
			err := database.DeleteFeed(db, i)
			if err != nil {
				logger.WithError(err).Info("Failed to delete feed")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error() + "\n"))
				return
			}
		}
		w.Write([]byte(http.StatusText(http.StatusOK) + "\n"))
	case "GET", "":
		messages, err := database.GetAllFeeds(db)
		if err != nil {
			logger.WithError(err).Info("Failed to get feeds")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error() + "\n"))
			return
		}
		for _, msg := range messages {
			w.Write([]byte(msg + "\n"))
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented) + "\n"))
	}
}

func handleShows(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	logger := utils.InitLogger()

	q := r.URL.Query()
	switch r.Method {
	case "POST":
		for _, i := range q["show"] {
			err := database.InsertShow(db, i)
			if err != nil {
				logger.WithError(err).Info("Failed to insert show")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error() + "\n"))
				return
			}
		}
		w.Write([]byte(http.StatusText(http.StatusOK) + "\n"))
	case "DELETE":
		for _, i := range q["show"] {
			err := database.DeleteShow(db, i)
			if err != nil {
				logger.WithError(err).Info("Failed to delete show")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error() + "\n"))
				return
			}
		}
		w.Write([]byte(http.StatusText(http.StatusOK) + "\n"))
	case "GET", "":
		messages, err := database.GetAllShows(db)
		if err != nil {
			logger.WithError(err).Info("Failed to get shows")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error() + "\n"))
			return
		}
		for _, msg := range messages {
			w.Write([]byte(msg + "\n"))
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented) + "\n"))
	}
}

func handleQuality(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	logger := utils.InitLogger()

	q := r.URL.Query()
	switch r.Method {
	case "POST":
		for _, i := range q["quality"] {
			qual, err := strconv.Atoi(i)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(http.StatusText(http.StatusBadRequest) + "\n"))
				return
			}
			err = database.InsertQuality(db, qual)
			if err != nil {
				logger.WithError(err).Info("Failed to insert quality")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error() + "\n"))
				return
			}
		}
		w.Write([]byte(http.StatusText(http.StatusOK) + "\n"))
	case "DELETE":
		for _, i := range q["quality"] {
			qual, err := strconv.Atoi(i)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(http.StatusText(http.StatusBadRequest) + "\n"))
				return
			}
			err = database.DeleteQuality(db, qual)
			if err != nil {
				logger.WithError(err).Info("Failed to delete quality")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error() + "\n"))
				return
			}
		}
		w.Write([]byte(http.StatusText(http.StatusOK) + "\n"))
	case "GET", "":
		messages, err := database.GetAllQualities(db)
		if err != nil {
			logger.WithError(err).Info("Failed to get quality")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error() + "\n"))
			return
		}
		for _, msg := range messages {
			w.Write([]byte(strconv.Itoa(msg) + "\n"))
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented) + "\n"))
	}
}

func handleHistory(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	logger := utils.InitLogger()

	q := r.URL.Query()
	switch r.Method {
	case "DELETE":
		var quality int
		var err error
		if q.Get("quality") != "" {
			quality, err = strconv.Atoi(q.Get("quality"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(http.StatusText(http.StatusBadRequest) + "\n"))
				return
			}
		}
		downloadInfo := scrapers.DownloadInfo{
			Show:     q.Get("show"),
			Episode:  q.Get("episode"),
			Quality:  quality,
			Uploader: q.Get("uploader"),
		}
		err = database.DeleteHistory(db, downloadInfo)
		if err != nil {
			logger.WithError(err).Info("Failed to delete history")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error() + "\n"))
			return
		}
		w.Write([]byte(http.StatusText(http.StatusOK) + "\n"))
	case "GET", "":
		var quality int
		var err error
		if q.Get("quality") != "" {
			quality, err = strconv.Atoi(q.Get("quality"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(http.StatusText(http.StatusBadRequest) + "\n"))
				return
			}
		}
		downloadInfo := scrapers.DownloadInfo{
			Show:     q.Get("show"),
			Episode:  q.Get("episode"),
			Quality:  quality,
			Uploader: q.Get("uploader"),
		}
		history, err := database.GetHistory(db, downloadInfo)
		if err != nil {
			logger.WithError(err).Info("Failed to get history")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error() + "\n"))
			return
		}
		json, err := json.Marshal(history)
		if err != nil {
			logger.WithError(err).Info("Failed to marshal history into JSON")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error() + "\n"))
			return
		}
		w.Write([]byte(string(json) + "\n"))
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented) + "\n"))
	}
}

func handleStart(w http.ResponseWriter, r *http.Request, runConfig config.Config) {
	var feeds []string
	var shows []string
	var quality []int

	logger := utils.InitLogger()

	switch r.Method {
	case "POST":
		q := r.URL.Query()
		if q.Get("feed") != "" {
			feeds = q["feed"]
		}
		if q.Get("show") != "" {
			shows = q["show"]
		}
		if q.Get("quality") != "" {
			for _, v := range q["quality"] {
				i, err := strconv.Atoi(v)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(http.StatusText(http.StatusBadRequest) + "\n"))
					return
				}
				quality = append(quality, i)
			}
		}
		go func() {
			if err := run.StartRun(runConfig, feeds, shows, quality, logger); err != nil {
				logger.WithError(err).Info("Run failed due to error")
			}
		}()
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(http.StatusText(http.StatusAccepted) + "\n"))
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented) + "\n"))
	}
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PONG\n"))
}
