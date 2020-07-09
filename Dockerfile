FROM alpine
WORKDIR /blog-crawler
COPY ./static ./static
COPY ./blogCrawlerConf.json ./blogCrawlerConf.json
COPY ./cmd ./cmd
ENV BLOG_CRAWLER_CONF ./blogCrawlerConf.json
EXPOSE 8003
CMD ["./cmd/blog-crawler"]