#!/bin/bash

go build -o blog-crawler *.go;

if [ -e ./blogCrawlerConf.json ]; then
  cp blog-crawler blogCrawlerConf.json $GOPATH/bin;
else
  cp blog-crawler $GOPATH/bin;
fi
