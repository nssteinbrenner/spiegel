package scrapers

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

const animeToshoBaseURL = "https://feed.animetosho.org/rss2"

func ScrapeAnimeTosho(c chan<- DownloadInfo, show string, quality []int) error {
	defer close(c)
	fp := newParser("animetosho.org")

	format := "%s?only_tor=1&q=%s&filter%%5B0%%5D%%5Bt%%5D=nyaa_class&filter%%5B0%%5D%%5Bv%%5D=trusted"
	url := fmt.Sprintf(format, animeToshoBaseURL, url.QueryEscape(show))
	feed, err := fp.ParseURL(url)
	if err != nil {
		return err
	}

	for _, item := range feed.Items {
		re := regexp.MustCompile("Torrent.*href=\"(magnet.*)\">Magnet")
		match := re.FindStringSubmatch(item.Description)
		if match == nil {
			continue
		}
		link := match[1]
		for _, v := range toDownloadInfo(item.Title) {
			for _, q := range quality {
				if strings.ToLower(v.Show) == strings.ToLower(show) && v.Quality >= q {
					v.Link = link
					c <- v
					break
				}
			}
		}
	}
	return nil
}
