package http

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nssteinbrenner/spiegel/config"
	"github.com/nssteinbrenner/spiegel/database"
	"github.com/nssteinbrenner/spiegel/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func HTTPServer(runConfig config.Config) error {
	db, err := database.Open(runConfig.DatabaseDirectory)
	if err != nil {
		return err
	}
	defer db.Close()

	mux := &http.ServeMux{}
	mux.HandleFunc("/feeds", func(w http.ResponseWriter, r *http.Request) {
		handleFeeds(w, r, db)
	})
	mux.HandleFunc("/shows", func(w http.ResponseWriter, r *http.Request) {
		handleShows(w, r, db)
	})
	mux.HandleFunc("/quality", func(w http.ResponseWriter, r *http.Request) {
		handleQuality(w, r, db)
	})
	mux.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		handleHistory(w, r, db)
	})
	mux.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		handleStart(w, r, runConfig)
	})
	mux.HandleFunc("/ping", handlePing)
	mux.Handle("/metrics", promhttp.Handler())

	var handler http.Handler = mux
	handler = metricsHandler(handler)

	if runConfig.HTTPEnabled {
		err := http.ListenAndServe(fmt.Sprintf(":%s", runConfig.HTTPPort), handler)
		if err != nil {
			return err
		}
	}

	if runConfig.HTTPSEnabled {
		_, err := os.Stat(runConfig.SSLCertificate)
		if err != nil {
			return err
		}
		_, err = os.Stat(runConfig.SSLCertificateKey)
		if err != nil {
			return err
		}
		err = http.ListenAndServeTLS(
			fmt.Sprintf(":%s", runConfig.HTTPSPort),
			runConfig.SSLCertificate,
			runConfig.SSLCertificateKey,
			handler,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func metricsHandler(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := statusWriter{ResponseWriter: w}
		handler.ServeHTTP(&sw, r)
		duration := time.Now().Sub(start).Seconds()
		if sw.status == 0 {
			sw.status = 200
		}
		trimmedURI := r.RequestURI
		if strings.Contains(r.RequestURI, "?") {
			trimmedURI = r.RequestURI[:strings.Index(r.RequestURI, "?")]
		}
		metrics.RequestDuration.WithLabelValues(
			trimmedURI,
			strconv.Itoa(sw.status),
			r.Method,
		).Observe(duration)
	}
}
