package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type strArrayFlags []string
type intArrayFlags []int

func SetArgs(config Config) (Config, bool, []string, []string, []int) {
	var feedArgs strArrayFlags
	var showArgs strArrayFlags
	var qualityArgs intArrayFlags

	oneshotRun := false

	flag.Usage = usage
	flag.StringVar(&config.DownloadDirectory, "downloaddir", config.DownloadDirectory, "")
	flag.StringVar(&config.TransmissionHost, "transmissionhost", config.TransmissionHost, "")
	flag.StringVar(&config.TransmissionPort, "transmissionport", config.TransmissionPort, "")
	flag.StringVar(&config.TransmissionUser, "transmissionuser", config.TransmissionUser, "")
	flag.StringVar(&config.TransmissionPassword, "transmissionpassword", config.TransmissionPassword, "")
	flag.BoolVar(&config.TransmissionHTTPS, "transmissionhttps", config.TransmissionHTTPS, "")
	flag.BoolVar(&config.HTTPEnabled, "httpserver", config.HTTPEnabled, "")
	flag.StringVar(&config.HTTPPort, "httpport", config.HTTPPort, "")
	flag.BoolVar(&config.HTTPSEnabled, "httpsserver", config.HTTPSEnabled, "")
	flag.StringVar(&config.HTTPSPort, "httpsport", config.HTTPSPort, "")
	flag.StringVar(&config.SSLCertificate, "sslcertificate", config.SSLCertificate, "")
	flag.StringVar(&config.SSLCertificateKey, "sslcertificatekey", config.SSLCertificateKey, "")

	oneshotCmd := flag.NewFlagSet("oneshot", flag.ExitOnError)
	oneshotCmd.Usage = usage
	oneshotCmd.Var(&feedArgs, "feed", "")
	oneshotCmd.Var(&showArgs, "show", "")
	oneshotCmd.Var(&qualityArgs, "quality", "")
	oneshotCmd.StringVar(&config.DownloadDirectory, "downloaddir", config.DownloadDirectory, "")
	oneshotCmd.StringVar(&config.TransmissionHost, "transmissionhost", config.TransmissionHost, "")
	oneshotCmd.StringVar(&config.TransmissionPort, "transmissionport", config.TransmissionPort, "")
	oneshotCmd.StringVar(&config.TransmissionUser, "transmissionuser", config.TransmissionUser, "")
	oneshotCmd.StringVar(&config.TransmissionPassword, "transmissionpassword", config.TransmissionPassword, "")
	oneshotCmd.BoolVar(&config.TransmissionHTTPS, "transmissionhttps", config.TransmissionHTTPS, "")

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "oneshot":
			oneshotCmd.Parse(os.Args[2:])
			oneshotRun = true
		default:
			flag.Parse()
		}
	}

	return config, oneshotRun, []string(feedArgs), []string(showArgs), []int(qualityArgs)
}

func usage() {
	fmt.Println(`Usage:
	spiegel [options]
	spiegel <subcommand> [options]

	subcommands:
		oneshot

	global options:
		-downloaddir				Path to directory to store downloads in
		-transmissionhost		   IP or hostname of Transmission server
		-transmissionport		   Port number to connect to on Transmission server
		-transmissionuser		   Username to login to Transmission
		-transmissionpassword	   Password for Transmission user
		-transmissionhttps		  Connect to Transmission over HTTPS

	non-subcommand options:
		-httpserver		 Listen on HTTP
		-httpport		   HTTP port
		-httpsserver		Listen on HTTPS
		-httpsport		  HTTPS port
		-sslcertificate	 Path to SSL certificate for HTTPS server
		-sslcertificatekey  Path to SSL certificate key for HTTPS server

	oneshot specific options:
		-feed	   Feeds to download from (can be specified multiple times)
		-show	   Shows to download (can be specified multiple times)
		-quality	Allowed video quality to download (can be specified multiple times)`)
}

func (i *strArrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *strArrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *intArrayFlags) Set(value string) error {
	res, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*i = append(*i, res)
	return nil
}

func (i *intArrayFlags) String() string {
	var res []string
	for _, val := range *i {
		conVal := strconv.Itoa(val)
		res = append(res, conVal)
	}
	return strings.Join(res[:], ",")
}
