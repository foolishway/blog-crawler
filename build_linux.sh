#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/blog-crawler *.go

cp ./cmd/blog-crawler ../blog-crawler-deploy/cmd