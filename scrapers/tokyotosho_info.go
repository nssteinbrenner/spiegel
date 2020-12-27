package scrapers

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

const tokyoToshoBaseURL = "https://www.tokyotosho.info/rss.php"

func ScrapeTokyoTosho(c chan<- DownloadInfo, show string, quality []int) error {
	defer close(c)
	fp := newParser("tokyotosho.info")

	format := "%s?filter=1&terms=%s%%20&searchComment=0&reversepolarity=1&entries=750"
	url := fmt.Sprintf(format, tokyoToshoBaseURL, url.QueryEscape(show))
	feed, err := fp.ParseURL(url)
	if err != nil {
		return err
	}

	for _, item := range feed.Items {
		found, err := regexp.Match(".*Authorized:\\sYes.*", []byte(item.Description))
		if err != nil {
			return err
		}
		if found {
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
	}
	return nil
}
