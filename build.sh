#!/bin/bash

go build -o blog-crawler *.go;

cp blog-crawler blogCrawlerConf.json $GOPATH/bin;