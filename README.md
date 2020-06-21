blog-crawler like its name, is a simple crawler, which can fetch blog title, address,  author, and publish time and so on  by your config file.

blog-crawler will ignore the blogs which has already got, so you can get the latest blog by execute the blog-crawler binary file.

### Usage

The key step is set the **BLOG_CRAWLER_CONF** environment variable，which points the config file, and then

you can execute the blog-crawler any where, if you did not set the **BLOG_CRAWLER_CONF** , blog-crawler will find the config file at the current path, if not found, it will log the error and then exit.

