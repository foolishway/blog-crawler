package crawler

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"os"

	//"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type blog struct {
	Author      string `json:"author"`
	Address     string `json:"address"`
	PageStyle   string `json:"pageStyle"`
	PageRule    string `json:pageRule`
	PostStyle   string `json:postStyle`
	TitleStyle  string `json:titleStyle`
	TimeStyle   string `json:timeStyle`
	PresentTime string `json:presentTime`
}
type Crawler struct {
	Blogs      []blog   `json:blogs`
	Exclude    []string `json:exclude`
	OutputType string   `outputType`
	Output     io.Writer
	Buf        *bytes.Buffer
	Mutex      sync.Mutex
}

var authorPresent map[string]time.Time

func (cr *Crawler) Start() {
	authorPresent = make(map[string]time.Time, len(cr.Blogs))
	blogs := cr.Blogs
	wg := &sync.WaitGroup{}
	for _, blog := range blogs {
		wg.Add(1)
		go cr.craw(&blog, 1, wg)
	}
	wg.Wait()
	//write cache buf to cache file
	fmt.Println("complete.")
	writeToCacheFile(cr.Buf)
}
func (cr *Crawler) craw(b *blog, pageNum int, wg *sync.WaitGroup) {
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
			if checkURL(b.PageRule) {
				addr = replacePageNum(b.PageRule, strconv.Itoa(pageNum))
			} else {
				u, _ := url.Parse(b.Address)
				addr = u.Host + replacePageNum(b.PageRule, strconv.Itoa(pageNum))
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
		panic(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	posts := doc.Find(b.PostStyle)

	if posts.Length() == 0 {
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
			return
		}

		timeStr := s.Find(b.TimeStyle).Text()

		address, _ := s.Find(b.TitleStyle).Attr("href")
		author := b.Author
		cr.writeToOutput(title, address, timeStr, author, cr.Output)
	})
	if !noNewBlog {
		wg.Add(1)
		initPageRule(b, doc)
		pageNum++
		cr.craw(b, pageNum, wg)
	}
}
func (cr *Crawler) writeToOutput(title, address, time string, author string, output io.Writer) {
	if output != nil {
		fmt.Fprintf(output, fmt.Sprintf("题目：%s；\n地址：%s；\n作者：%s；\n发布时间：%s\n\n", title, address, author, time))
	}
	cr.Mutex.Lock()
	cr.Buf.Write([]byte(cacheFormat(title, author)))
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

func initPageRule(b *blog, doc *goquery.Document) {
	if b.PageRule == "" {
		if ps := doc.Find(b.PageStyle); ps.Length() > 1 {
			attr := ps.Get(1).Attr
			for _, node := range attr {
				if node.Key == "href" {
					b.PageRule = node.Val
				}
			}
		}
	}
}
func writeToCacheFile(buf *bytes.Buffer) {
	f, err := os.OpenFile("./blog.cache", os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		panic(fmt.Sprintf("Write to cache file error: %v", err))
	}
	defer f.Close()
	_, err = f.Write(buf.Bytes())
	if err != nil {
		panic(fmt.Sprintf("Write to cache file error: %v", err))
	}
}
