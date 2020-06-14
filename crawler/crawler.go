package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/url"
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
}

var authorPresent map[string]time.Time

func (cr *Crawler) Start() {
	authorPresent = make(map[string]time.Time, len(cr.Blogs))
	blogs := cr.Blogs
	wg := &sync.WaitGroup{}
	for _, blog := range blogs {
		wg.Add(1)
		go craw(cr.Exclude, &blog, 1, cr.Output, wg)
	}
	wg.Wait()
	fmt.Println("complete.")
}

func craw(exs []string, b *blog, pageNum int, output io.Writer, wg *sync.WaitGroup) {
	defer wg.Done()
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
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
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
		if checkExclude(title, exs) {
			return
		}

		timeStr := s.Find(b.TimeStyle).Text()
		if timeStr == "" {
			//TODO get first blog
			noNewBlog = true
			return
		}

		time, err := formatTime(timeStr)
		if err != nil {
			//TODO get first blog
			noNewBlog = true
			return
		}
		pTime, err := formatTime(b.PresentTime)
		if err != nil {
			panic("Convert present time error.")
		}
		if time.Before(pTime) {
			noNewBlog = true
			return
		}

		//record the present time
		if authorPresent[b.Author].IsZero() {
			authorPresent[b.Author] = time
		}

		address, _ := s.Find(b.TitleStyle).Attr("href")
		author := b.Author
		writeToOutput(title, address, timeStr, author, output)
	})
	if !noNewBlog {
		wg.Add(1)
		initPageRule(b, doc)
		pageNum++
		craw(exs, b, pageNum, output, wg)
	}
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

func getNum(path string) (int, error) {
	numRg := regexp.MustCompile(`(\d+)`)
	if !numRg.MatchString(path) {
		return 0, fmt.Errorf("Compile path error.")
	}
	num := numRg.FindAllString(path, -1)
	if len(num) >= 1 {
		n, _ := strconv.Atoi(num[0])
		return n, nil
	}
	return 0, fmt.Errorf("Compile path error.")
}
func checkNum(pageNum string) bool {
	numRg := regexp.MustCompile(`[1-9]\d+|^[0-9]$`)
	isNum := numRg.MatchString(pageNum)
	return isNum
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
func formatTime(timeStr string) (time.Time, error) {
	exp, err := regexp.Compile(`^\d{4}([年月日/.]{1})\d{1,2}([年月日/.]{1})\d{1,2}([年月日/.])?$`)
	if err != nil || !exp.MatchString(timeStr) {
		return time.Now(), err
	}
	rep := []string{"年", "月", "日", "/", "."}

	for _, r := range rep {
		timeStr = strings.Replace(timeStr, r, "-", -1)
	}
	if timeStr[len(timeStr)-1:] == "-" {
		timeStr = timeStr[:len(timeStr)-1]
	}
	t, err := time.Parse("2006-1-2", timeStr)

	if err != nil {
		return time.Now(), err
	}
	return t, nil
}
func writeToOutput(title, address, time string, author string, output io.Writer) {
	if output != nil {
		fmt.Fprintf(output, fmt.Sprintf("题目：%s；\n地址：%s；\n作者：%s；\n发布时间：%s\n\n", title, address, author, time))
	}
}
