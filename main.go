package main

import (
	"blog-crawler/crawler"
	"blog-crawler/models"

	"encoding/json"
	"io/ioutil"
	"log"
	"path"

	//"blog-crawler/crawler"
	//"bytes"
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
	//"log"
	"os"
	//"path"
)

var (
	cachePath string = "./blog.cache"
	confPath  string = "./blogCrawlerConf.json"
)

const avg = 10

func main() {
	//s := &hs.HttpServer{"8003"}
	//s.StartServer()
	envConf, envSet := os.LookupEnv("BLOG_CRAWLER_CONF")
	//if set BLOG_CRAWER_CONF environment variable, the cache file will be generated under the same path
	if envSet {
		confPath = envConf
		cachePath = path.Dir(confPath) + "/" + path.Base(cachePath)
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

	c := &crawler.Crawler{CachePath: cachePath, CollectArticles: make([]models.Article, 0)}
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

	//if c.Buf == nil {
	//	var f *os.File
	//	if !fileExists(cachePath) {
	//		f, err = os.Create(cachePath)
	//		if err != nil {
	//			panic("create cache file error.")
	//		}
	//	} else {
	//		f, err = os.Open(cachePath)
	//		if err != nil {
	//			panic("open cache file error.")
	//		}
	//	}
	//
	//	defer f.Close()
	//	cacheBytes, err := ioutil.ReadAll(f)
	//	if err != nil {
	//		panic(fmt.Sprintf("Read blog.cache file error: %v", err))
	//	}
	//	c.Buf = bytes.NewBuffer(cacheBytes)
	//}

	c.Start()
}
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
