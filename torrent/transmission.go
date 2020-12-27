package torrent

import (
	"strconv"

	"github.com/hekmon/transmissionrpc"
)

func GetTransmissionConnection(host, port, user, password string, https bool) (*transmissionrpc.Client, error) {
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	transmissionbt, err := transmissionrpc.New(host, user, password, &transmissionrpc.AdvancedConfig{
		HTTPS: https,
		Port:  uint16(intPort),
	})
	if err != nil {
		return nil, err
	}
	return transmissionbt, nil
}

func AddTransmissionTorrent(tbt *transmissionrpc.Client, downloadDirectory, link string) error {
	_, err := tbt.TorrentAdd(&transmissionrpc.TorrentAddPayload{
		Filename:    &link,
		DownloadDir: &downloadDirectory,
	})
	if err != nil {
		return err
	}
	return nil
}
