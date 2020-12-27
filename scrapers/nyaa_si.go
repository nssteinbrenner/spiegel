package scrapers

import (
	"fmt"
	"net/url"
	"strings"
)

const nyaaBaseURL = "https://nyaa.si/?page=rss"

func ScrapeNyaa(c chan<- DownloadInfo, show string, quality []int) error {
	defer close(c)
	fp := newParser("nyaa.si")

	url := fmt.Sprintf("%s&q=%s&c=1_2&f=2", nyaaBaseURL, url.QueryEscape(show))
	feed, err := fp.ParseURL(url)
	if err != nil {
		return err
	}

	for _, item := range feed.Items {
		for _, v := range toDownloadInfo(item.Title) {
			for _, q := range quality {
				if strings.ToLower(v.Show) == strings.ToLower(show) && v.Quality >= q {
					v.Link = item.Link
					c <- v
					break
				}
			}
		}
	}
	return nil
}
