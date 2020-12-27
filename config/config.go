package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DownloadDirectory    string `envconfig:"DOWNLOAD_DIRECTORY"`
	DatabaseDirectory    string `default:"" envconfig:"DATABASE_DIRECTORY"`
	TransmissionHost     string `default:"127.0.0.1" envconfig:"TRANSMISSION_HOST"`
	TransmissionPort     string `default:"9091" envconfig:"TRANSMISSION_PORT"`
	TransmissionUser     string `default:"" envconfig:"TRANSMISSION_USER"`
	TransmissionPassword string `default:"" envconfig:"TRANSMISSION_PASSWORD"`
	TransmissionHTTPS    bool   `default:false envconfig:"TRANSMISSION_HTTPS"`
	HTTPEnabled          bool   `default:false envconfig:"HTTPENABLED"`
	HTTPPort             string `default:"80" envconfig:"HTTPPORT"`
	HTTPSEnabled         bool   `default:false envconfig:"HTTPSENABLED"`
	HTTPSPort            string `default:"443" envconfig:"HTTPSPORT"`
	SSLCertificate       string `default:"" envconfig:"SSLCERTIFICATE"`
	SSLCertificateKey    string `default:"" envconfig:"SSLCERTIFICATEKEY"`
}

func BuildConfig() (Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
