FROM alpine
WORKDIR /blog-crawler
COPY ./static ./static
COPY ./blogCrawlerConf.json ./blogCrawlerConf.json
COPY ./cmd ./cmd


EXPOSE 8003
CMD ["./cmd/blog-crawler"]