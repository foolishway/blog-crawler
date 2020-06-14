package main

import (
	"blog-crawler/crawler"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	conf, err := os.Open("./conf.json")
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
	c.Start()
}
