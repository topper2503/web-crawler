package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Fetcher interface {
	Fetch(*url.URL) ([]string, error)
}

type FetcherImpl struct {
	client *http.Client
}

func NewFetcherImpl() Fetcher {
	return FetcherImpl{
		client: http.DefaultClient,
	}
}

func (f FetcherImpl) Fetch(url *url.URL) ([]string, error) {
	resp, err := f.client.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non 200 status code for %s: %d", url, resp.StatusCode)
	}

	return parsePageForLinks(resp.Body)
}

func parsePageForLinks(data io.Reader) ([]string, error) {
	links := []string{}
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		return []string{}, err
	}

	doc.Find("a[href]").Each(func(i int, item *goquery.Selection) {
		href, _ := item.Attr("href")

		links = append(links, href)
	})
	return links, nil
}
