FROM alpine
WORKDIR /blog-crawler
COPY ./static ./static
COPY ./blogCrawlerConf.json ./blogCrawlerConf.json
COPY ./cmd ./cmd
ENV BLOG_CRAWLER_CONF ./blogCrawlerConf.json
ENV ROBOT_BASE_PATH ""
ENV ROBOT_ACCESS_TOKEN ""
ENV ROBOT_ACCESS_KEY ""
ENV SESSION_KEY ""

EXPOSE 8003
CMD ["./cmd/blog-crawler"]