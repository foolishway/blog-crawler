blog-crawler 通过给定的配置文件，爬取指定的博客文章。每次爬取都是博主最新的文章。类似于"订阅"。

**_tag 1.0.0_** 是命令行版本，命令行下运行blog-crawler，会把抓取到的文章信息输出到stdOut.

用法：

配置环境变量：_**BLOG_CRAWLER_CONF**_ 指向一个配置文件，配置文件格式如下：
``` json
{
  "blogs": [{
      "address": "http://www.tracefact.net/",
      "author": "张子阳",
      "pageRule": "?page=1",
      "postStyle": "div.article",
      "titleStyle": "div.article .title a"
    },
    {
      "address": "https://www.zhangxinxu.com/wordpress/",
      "author": "张鑫旭",
      "pageRule": "wordpress/page/2/",
      "postStyle": "div.post",
      "titleStyle": "div.post h2 a",
      "timeStyle": "div.post .date"
    }]
}
```
字段说明：
* address 博客地址
* author 作者
* pageRule 分页规则
* postStyle 文章列表样式
* titleStyle 博客标题样式
* timeStyle 博客发表时间样式

