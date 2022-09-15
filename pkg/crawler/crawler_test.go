package crawler_test

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/topper2503/web-crawler/pkg/crawler"
	"github.com/topper2503/web-crawler/pkg/logging"
)

func TestMain(m *testing.M) {
	logging.SetupLogger()
	exitVal := m.Run()
	os.Exit(exitVal)
}

type FakeFetcher struct {
	urls map[string][]string
}

func (f FakeFetcher) Fetch(url *url.URL) ([]string, error) {
	val, ok := f.urls[url.String()]
	if !ok {
		return []string{}, fmt.Errorf("not found")
	}
	return val, nil
}

func TestCrawler_Run(t *testing.T) {
	td := []struct {
		name               string
		url                string
		fetcher            FakeFetcher
		depth              int
		expectedTotalLinks int
	}{
		{
			name: "Follows correct depth",
			url:  "https://ThisIsARealUrl",
			fetcher: FakeFetcher{urls: map[string][]string{
				"https://ThisIsARealUrl":   {"/a", "/b"},
				"https://ThisIsARealUrl/a": {"/b"},
				"https://ThisIsARealUrl/b": {"/c", "/e"},
			},
			},
			depth:              2,
			expectedTotalLinks: 3,
		},
		{
			name: "Follows extended depth",
			url:  "https://ThisIsARealUrl",
			fetcher: FakeFetcher{urls: map[string][]string{
				"https://ThisIsARealUrl":   {"/a", "/b"},
				"https://ThisIsARealUrl/a": {"/b"},
				"https://ThisIsARealUrl/b": {"/c", "/e"},
			},
			},
			depth:              3,
			expectedTotalLinks: 5,
		},
	}

	for _, tc := range td {
		t.Run(tc.name, func(*testing.T) {
			u, err := url.Parse(tc.url)
			assert.NoError(t, err)

			c := crawler.New(tc.fetcher, u, tc.depth)

			page, err := c.Run()
			assert.NoError(t, err)
			page.PrintPages()

			assert.EqualValues(t, tc.expectedTotalLinks, getTotalPageCount(page))
		})
	}
}

func getTotalPageCount(p *crawler.Page) int {
	count := 1

	for _, page := range p.Links {
		count = count + getTotalPageCount(page)
	}
	return count
}
