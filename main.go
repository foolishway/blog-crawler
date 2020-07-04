package main

import (
	"blog-crawler/crawler"
	hs "blog-crawler/http-server"
	"blog-crawler/models"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	confPath string = "./blogCrawlerConf.json"
)

func main() {
	envConf, envSet := os.LookupEnv("BLOG_CRAWLER_CONF")
	//if set BLOG_CRAWER_CONF environment variable, the cache file will be generated under the same path
	if envSet {
		confPath = envConf
	}
	_, err := os.Stat(confPath)
	if err != nil && os.IsNotExist(err) {
		log.Fatalf("Config file not found.")
	}

	//init crawle
	go startCrawl()

	tick := time.Tick(24 * time.Hour)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			select {
			case <-tick:
				go startCrawl()
			case <-signalChan:
				models.CloseDb()
				os.Exit(0)
			}
		}
	}()

	defer models.CloseDb()

	//start http server
	s := &hs.HttpServer{"8003"}
	s.StartServer()
}
func startCrawl() {
	log.Println("Start crawl...")
	c := crawler.NewCrawler(confPath)
	c.Start()
}
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
