package main

import (
	"blog-crawler/crawler"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	confPath := "./blogCrawlerConf.json"
	envConf, envSet := os.LookupEnv("BLOG_CRAWER_CONF")
	if envSet {
		confPath = envConf
	}
	_, err := os.Stat(confPath)
	if err != nil && os.IsNotExist(err) {
		log.Fatalf("Config file not found.")
	}

	conf, err := os.Open(confPath)
	if err != nil {
		log.Fatalf("Open conf error %v", err)
	}
	defer conf.Close()

	c := &crawler.Crawler{}
	b, err := ioutil.ReadAll(conf)
	if err != nil {
		log.Fatalf("Read from conf error %v", err)
	}
	err = json.Unmarshal(b, c)
	if err != nil {
		log.Fatalf("Unmarshall json error: %v", err)
	}

	if c.OutputType == "stdout" {
		c.Output = os.Stdout
	}

	if c.Buf == nil {
		cachePath := "./blog.cache"
		var f *os.File
		if !fileExists(cachePath) {
			f, err = os.Create("./blog.cache")
			if err != nil {
				panic("create cache file error.")
			}
		} else {
			f, err = os.Open(cachePath)
			if err != nil {
				panic("open cache file error.")
			}
		}

		defer f.Close()
		cacheBytes, err := ioutil.ReadAll(f)
		if err != nil {
			panic(fmt.Sprintf("Read blog.cache file error: %v", err))
		}
		c.Buf = bytes.NewBuffer(cacheBytes)
	}

	time.Tick(24 * time.Hour)
	c.Start()
}
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
