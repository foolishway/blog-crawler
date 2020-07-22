package main

import (
	"blog-crawler/crawler"
	"blog-crawler/duty"
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
	initLog()
	envConf, envSet := os.LookupEnv("BLOG_CRAWLER_CONF")
	//if set BLOG_CRAWER_CONF environment variable, the cache file will be generated under the same path
	if envSet {
		confPath = envConf
	}
	_, err := os.Stat(confPath)
	if err != nil && os.IsNotExist(err) {
		log.Fatalf("Config file not found.")
	}

	//init crawl
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

	//start duty notify
	go duty.StartDuty()

	//start http server
	s := &hs.HttpServer{"8003"}
	s.StartServer()
}

func startCrawl() {
	log.Println("Start crawl...")
	c := crawler.NewCrawler(confPath)
	c.Start()
}

func initLog() {
	logFile, err := os.OpenFile("./blog-crawler.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	defer logFile.Close()

	if err != nil {
		log.Printf("init log error %v\n", err)
		return
	}
	log.SetOutput(logFile)
}
