package config

import (
	"flag"

	"github.com/topper2503/web-crawler/pkg/env"
	"github.com/topper2503/web-crawler/pkg/logging"

	"go.uber.org/zap"
)

type Config struct {
	TargetURL  string
	CrawlDepth int
}

var (
	config Config
)

func Initialise() {
	logger := logging.GetLogger()

	var (
		targetURL  string
		crawlDepth int
	)

	flag.StringVar(&targetURL, "targetURL", env.GetenvString("URL", ""), "URL that the crawl will be started on")
	flag.IntVar(&crawlDepth, "crawlDepth", env.GetenvInt("DEPTH", 2), "How deep the recursive crawler should search")

	logger.Info("configuration",
		zap.String("targetURL", targetURL),
		zap.Int("crawlDepth", crawlDepth),
	)

	config.TargetURL = targetURL
	config.CrawlDepth = crawlDepth
}

func Get() Config {
	return config
}
