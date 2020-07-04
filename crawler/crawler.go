package crawler

import (
	"blog-crawler/models"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const avg = 10

var existMap map[string]struct{}

func fillExistMap() {
	articles := models.GetAllArticles()
	existMap = models.AriticleModelToMap(articles)
}

type Crawler struct {
	Blogs           []models.Blog `json:blogs`
	Exclude         []string      `json:exclude`
	OutputType      string        `outputType`
	Output          io.Writer
	Mutex           sync.Mutex
	CollectArticles []models.Article
}

func (cr *Crawler) Start() {
	fillExistMap()
	blogs := cr.Blogs
	wgQue := make([]*sync.WaitGroup, 0)
	for i := 0; i < len(blogs); i++ {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		wgQue = append(wgQue, wg)
		go cr.craw(wg, &(blogs[i]), 1)
	}
	//wait all blogs complete
	for i := 0; i < len(wgQue); i++ {
		wgQue[i].Wait()
	}
	fmt.Println("Crawl complete.")
	//writeToCacheFile(cr.Buf, cr.CachePath)
	if len(cr.CollectArticles) > 0 {
		//fmt.Println(cr.CollectArticles)
		models.InsertCollectArticles(cr.CollectArticles)
	}
}
func (cr *Crawler) craw(wg *sync.WaitGroup, b *models.Blog, pageNum int) {
	//fmt.Println("craw len(cr.CollectArticles)", len(cr.CollectArticles))
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
		wg.Done()
	}()
	var addr string
	if pageNum == 1 {
		addr = b.Address
	} else {
		//rule = "http://www....."、 "?pn=20"
		if b.PageRule != "" {
			if b.PageRule[0] != '/' {
				b.PageRule = "/" + b.PageRule
			}
			if checkURL(b.PageRule) {
				addr = replacePageNum(b.PageRule, strconv.Itoa(pageNum))
			} else {
				u, _ := url.Parse(b.Address)
				addr = u.Scheme + "://" + u.Host + replacePageNum(b.PageRule, strconv.Itoa(pageNum))
			}
		} else {
			return
		}
	}

	var stopCrawl bool
	// Request the HTML page.
	res, err := http.Get(addr)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		panic(fmt.Sprintf("Requst %s status code error: %d %s", b.Address, res.StatusCode, res.Status))
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	posts := doc.Find(b.PostStyle)

	if posts.Length() == 0 {
		if pageNum == 1 {
			log.Printf("%s dit not published articles yet.")
		}
		stopCrawl = true
		return
	}
	// Find blogs
	doc.Find(b.PostStyle).Each(func(i int, s *goquery.Selection) {
		if stopCrawl {
			return
		}
		title := s.Find(b.TitleStyle).Text()
		author := b.Author
		//filter exclude
		if checkExclude(title, cr.Exclude) {
			return
		}
		//report whether cache file contains the blog
		if isExist(title+"_"+b.Author) || overArg(author, cr.CollectArticles) {
			stopCrawl = true
			return
		}

		timeStr := s.Find(b.TimeStyle).Text()

		address, _ := s.Find(b.TitleStyle).Attr("href")
		if address != "" && !checkURL(address) {
			if address[0] != '/' {
				address = "/" + address
			}
			u, err := url.Parse(b.Address)
			if err != nil {
				panic("parse blog address error.")
			}
			address = u.Scheme + "://" + u.Host + address
		}
		cr.CollectArticles = append(cr.CollectArticles, models.Article{Title: title, Author: author, Address: address, PublishTime: timeStr})
	})
	if !stopCrawl {
		wg.Add(1)
		//initPageRule(b, doc)
		pageNum++
		cr.craw(wg, b, pageNum)
	}
}

func NewCrawler(confPath string) *Crawler {
	conf, err := os.Open(confPath)
	if err != nil {
		log.Fatalf("Open conf error %v", err)
	}
	defer conf.Close()

	c := &Crawler{CollectArticles: make([]models.Article, 0)}
	b, err := ioutil.ReadAll(conf)
	if err != nil {
		log.Fatalf("Read from conf error %v", err)
	}
	err = json.Unmarshal(b, c)
	if err != nil {
		log.Fatalf("Unmarshall json error: %v", err)
	}
	return c
}
func checkExclude(exStr string, exs []string) bool {
	for _, ex := range exs {
		if exStr == ex {
			return true
		}
	}
	return false
}
func checkURL(path string) bool {
	urlRg, _ := regexp.Compile(`^(?i)https?://.+`)
	return urlRg.MatchString(path)
}
func replacePageNum(uri, newPage string) string {
	rg := regexp.MustCompile(`\d+`)
	ms := rg.FindStringSubmatch(uri)
	//if len(ms) != 0 {
	//	ms = ms[1:]
	//}
	for _, m := range ms {
		uri = strings.Replace(uri, m, newPage, -1)
	}
	return uri
}
func isExist(key string) bool {
	_, exist := existMap[key]
	return exist
}
func overArg(author string, cellectAticle []models.Article) bool {
	var count int32
	for _, article := range cellectAticle {
		if article.Author == author {
			count++
		}
	}
	if count >= avg {
		return true
	}
	return false
}
