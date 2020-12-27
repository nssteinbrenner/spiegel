package scrapers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/nssteinbrenner/anitogo"
	"github.com/nssteinbrenner/spiegel/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; rv:78.0) Gecko/20100101 Firefox/78.0"

type DownloadInfo struct {
	Link     string
	Show     string
	Episode  string
	Quality  int
	Uploader string
	Batch    bool
}

type UserAgentTransport struct {
	http.RoundTripper
}

func (c *UserAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", userAgent)
	return c.RoundTripper.RoundTrip(r)
}

func getHTTPClient(scraper string) *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: promhttp.InstrumentRoundTripperDuration(
			metrics.ScraperRequestDuration.MustCurryWith(prometheus.Labels{"scraper": scraper}),
			&UserAgentTransport{http.DefaultTransport},
		),
	}
}

func newParser(scraper string) *gofeed.Parser {
	fp := gofeed.NewParser()
	fp.Client = getHTTPClient(scraper)

	return fp
}

func toDownloadInfo(filename string) []DownloadInfo {
	parsed := anitogo.Parse(filename, anitogo.DefaultOptions)
	if len(parsed.EpisodeNumber) == 0 && len(parsed.EpisodeNumberAlt) == 0 {
		return nil
	}

	quality := "0"
	re := regexp.MustCompile("^\\d{3,4}x(\\d{3,4})[pP]?$")
	match := re.FindStringSubmatch(parsed.VideoResolution)
	if match != nil {
		quality = match[1]
	}
	re = regexp.MustCompile("^(\\d{3,4})[pP]?$")
	match = re.FindStringSubmatch(parsed.VideoResolution)
	if match != nil {
		quality = match[1]
	}
	intQuality, err := strconv.Atoi(quality)
	if err != nil {
		return nil
	}

	episodeNumbers := parsed.EpisodeNumberAlt
	if len(parsed.EpisodeNumber) >= 1 && len(parsed.EpisodeNumberAlt) == 0 {
		episodeNumbers = parsed.EpisodeNumber
	}
	animeTitle := parsed.AnimeTitle
	for _, v := range parsed.AnimeSeason {
		animeTitle = fmt.Sprintf("%s S%s", animeTitle, v)
	}

	batch := false
	if len(episodeNumbers) > 1 {
		batch = true
	}

	ret := make([]DownloadInfo, len(episodeNumbers))
	if len(episodeNumbers) > 1 {
		intStart, startErr := strconv.Atoi(episodeNumbers[0])
		intEnd, endErr := strconv.Atoi(episodeNumbers[len(episodeNumbers)-1])
		if startErr == nil && endErr == nil {
			for i := intStart; i <= intEnd; i++ {
				ret = append(ret, DownloadInfo{
					Show:     animeTitle,
					Episode:  strconv.Itoa(i),
					Quality:  intQuality,
					Uploader: parsed.ReleaseGroup,
					Batch:    batch,
				})
			}
			return ret
		}
	}

	for _, episodeNumber := range episodeNumbers {
		episodeNumber = strings.TrimLeft(episodeNumber, "0")
		ret = append(ret, DownloadInfo{
			Show:     animeTitle,
			Episode:  episodeNumber,
			Quality:  intQuality,
			Uploader: parsed.ReleaseGroup,
			Batch:    batch,
		})
	}

	return ret
}
