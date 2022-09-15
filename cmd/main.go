package main

import (
	"log"
	"net/url"
	"os"

	"github.com/topper2503/web-crawler/cmd/config"
	"github.com/topper2503/web-crawler/pkg/crawler"
	"github.com/topper2503/web-crawler/pkg/fetcher"
	"github.com/topper2503/web-crawler/pkg/logging"
	"go.uber.org/zap"
)

func main() {
	logger := setup()
	logger.Info("running crawler")
	if config.Get().TargetURL == "" {
		logger.Fatal("no targetURL set")
		return
	}

	targetURL, err := url.Parse(config.Get().TargetURL)
	if err != nil {
		logger.Fatal("couldn't parse targetURL")
		return
	}

	webCrawler := crawler.New(fetcher.NewFetcherImpl(), targetURL, config.Get().CrawlDepth)
	page, err := webCrawler.Run()
	if err != nil {
		logger.Fatal("failed to run crawler", zap.Error(err))
	}
	page.PrintPages()
}

func setup() *zap.Logger {
	logger, err := logging.NewLogger(os.Getenv("ENV"))
	if err != nil {
		log.Fatalf("Cannot set up logger: %s", err.Error())
	}
	logging.SetLogger(logger)

	config.Initialise()

	return logger
}
