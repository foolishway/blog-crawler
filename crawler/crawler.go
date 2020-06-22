package crawler

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	//"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type blog struct {
	Author     string `json:"author"`
	Address    string `json:"address"`
	PageRule   string `json:pageRule`
	PostStyle  string `json:postStyle`
	TitleStyle string `json:titleStyle`
	TimeStyle  string `json:timeStyle`
	wg         *sync.WaitGroup
}
type Crawler struct {
	Blogs      []blog   `json:blogs`
	Exclude    []string `json:exclude`
	OutputType string   `outputType`
	Output     io.Writer
	Buf        *bytes.Buffer
	Mutex      sync.Mutex
	CachePath  string
}

func (cr *Crawler) Start() {
	blogs := cr.Blogs
	for i := 0; i < len(blogs); i++ {
		blogs[i].wg = &sync.WaitGroup{}
		blogs[i].wg.Add(1)
		go cr.craw(&(blogs[i]), 1)
	}
	//wait all blogs complete
	for i := 0; i < len(blogs); i++ {
		blogs[i].wg.Wait()
	}
	fmt.Println("complete.")
	writeToCacheFile(cr.Buf, cr.CachePath)
}
func (cr *Crawler) craw(b *blog, pageNum int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
		b.wg.Done()
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

	var noNewBlog bool
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
		noNewBlog = true
		return
	}
	// Find blogs
	doc.Find(b.PostStyle).Each(func(i int, s *goquery.Selection) {
		title := s.Find(b.TitleStyle).Text()
		//filter exclude
		if checkExclude(title, cr.Exclude) {
			return
		}
		//report whether cache file contains the blog
		if bytes.Contains(cr.Buf.Bytes(), []byte(cacheFormat(title, b.Author))) {
			noNewBlog = true
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
		author := b.Author
		cr.writeToOutput(title, address, timeStr, author, cr.Output)
	})
	if !noNewBlog {
		b.wg.Add(1)
		//initPageRule(b, doc)
		pageNum++
		cr.craw(b, pageNum)
	}
}
func (cr *Crawler) writeToOutput(title, address, time string, author string, output io.Writer) {
	if title == "" || address == "" {
		return
	}
	if output != nil {
		t := title
		ad := address
		au := author
		ti := time
		if title == "" {
			t = "--"
		}
		if address == "" {
			ad = "--"
		}
		if author == "" {
			au = "--"
		}
		if time == "" {
			ti = "--"
		}
		fmt.Fprintf(output, fmt.Sprintf("题目：%s；\n地址：%s；\n作者：%s；\n发布时间：%s\n\n", t, ad, au, ti))
	}
	cr.Mutex.Lock()
	c := cacheFormat(title, author)
	if !bytes.Contains(cr.Buf.Bytes(), []byte(c)) {
		cr.Buf.Write([]byte(c))
	}
	cr.Mutex.Unlock()
}

func cacheFormat(title, author string) string {
	return fmt.Sprintf("[[%s_%s]]", title, author)
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
func writeToCacheFile(buf *bytes.Buffer, cachePath string) {
	f, err := os.OpenFile(cachePath, os.O_TRUNC|os.O_WRONLY, 0644)

	if err != nil {
		panic(fmt.Sprintf("Write to cache file error: %v", err))
	}
	defer f.Close()
	_, err = f.Write(buf.Bytes())
	if err != nil {
		panic(fmt.Sprintf("Write to cache file error: %v", err))
	}
}
