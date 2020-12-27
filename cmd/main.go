package main

import (
	"os"
	"syscall"

	"github.com/nssteinbrenner/spiegel/config"
	"github.com/nssteinbrenner/spiegel/http"
	"github.com/nssteinbrenner/spiegel/run"
	"github.com/nssteinbrenner/spiegel/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	syscall.Umask(0)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	runConfig, err := config.BuildConfig()
	if err != nil {
		logrus.WithError(err).Info("Failed to get config")
		return
	}
	runConfig, oneshot, feeds, shows, quality := config.SetArgs(runConfig)
	if runConfig.HTTPEnabled || runConfig.HTTPSEnabled {
		if err := http.HTTPServer(runConfig); err != nil {
			logrus.WithError(err).Info("Failed to start HTTP server")
		}
	} else if oneshot {
		logger := utils.InitLogger()
		if err := run.StartRun(runConfig, feeds, shows, quality, logger); err != nil {
			logger.WithError(err).Info("Run failed due to error")
		}
	}
}
