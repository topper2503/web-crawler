package crawler

import (
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/topper2503/web-crawler/pkg/fetcher"
	"github.com/topper2503/web-crawler/pkg/logging"
	"go.uber.org/zap"
)

type Crawler struct {
	Page *Page

	seenUrls SeenURLs
	fetcher  fetcher.Fetcher
	depth    int
}

type Page struct {
	URL   *url.URL
	Links []*Page
}

type SeenURLs struct {
	List map[string]struct{} //valueless map, for checking if URL has already been seen
	m    sync.RWMutex        //for threadsafe read and write access to the list
}

func (s *SeenURLs) Exists(val string) bool {
	s.m.RLock()
	defer s.m.RUnlock()
	_, exists := s.List[val]
	return exists
}

func (s *SeenURLs) Set(val string) {
	s.m.Lock()
	defer s.m.Unlock()
	s.List[val] = struct{}{}
}

func New(fetcher fetcher.Fetcher, target *url.URL, depth int) *Crawler {
	return &Crawler{
		seenUrls: SeenURLs{List: map[string]struct{}{}},
		Page:     &Page{URL: target},
		fetcher:  fetcher,
		depth:    depth,
	}
}

func (c *Crawler) Run() (*Page, error) {
	start := time.Now()
	logger := logging.GetLogger()
	c.seenUrls.Set(c.Page.URL.String())
	err := c.crawlPage(c.Page, c.depth)
	if err != nil {
		return nil, err
	}

	elapsed := time.Since(start)
	logger.Info("Unique links crawled:", zap.Int("unique Links", len(c.seenUrls.List)))
	logger.Info("Crawling took", zap.Duration("elapsedTime", elapsed))
	return c.Page, nil
}

func (p *Page) PrintPages() {
	printPage(p, 0)
}

func printPage(page *Page, indent int) {
	logger := logging.GetLogger()

	a := strings.Join([]string{strings.Repeat("    ", indent), page.URL.String()}, "")
	logger.Info(a)

	for _, subpage := range page.Links {
		printPage(subpage, indent+1)
	}
}

func (c *Crawler) crawlPage(currentPage *Page, depth int) error {
	logger := logging.GetLogger()
	var wg sync.WaitGroup

	if depth <= 0 {
		return nil
	}

	pageLinks, err := c.fetcher.Fetch(currentPage.URL)
	if err != nil {
		return err
	}

	linksChan := make(chan *Page)
	var linksWg sync.WaitGroup

	linksWg.Add(1)
	defer linksWg.Done()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for link := range linksChan {
			currentPage.Links = append(currentPage.Links, link)
		}
	}()

	for _, link := range pageLinks {
		relURL, err := url.Parse(link)
		if err != nil {
			logger.Error("failed to parse URL on page", zap.String("URL", link), zap.String("page", currentPage.URL.String()), zap.Error(err))
			continue
		}
		newPageURl := currentPage.URL.ResolveReference(relURL)

		if c.seenUrls.Exists(newPageURl.String()) {
			continue
		}

		newPage := &Page{URL: newPageURl}
		c.crawlPage(newPage, depth-1)
		linksChan <- newPage
		c.seenUrls.Set(newPageURl.String())
	}
	return nil
}
