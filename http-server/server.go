package http_server

import (
	"blog-crawler/models"
	"blog-crawler/robot"
	"blog-crawler/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

type HttpServer struct {
	Port string
}

func (hs *HttpServer) StartServer() {
	//static server
	serveFile()
	tpl := template.Must(
		template.ParseGlob(
			"./static/views/*.html",
		),
	)

	//index
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexHandler(w, r, tpl)
	})
	//share
	http.HandleFunc("/share", shareHandler)
	err := http.ListenAndServe(":"+hs.Port, nil)
	if err != nil {
		log.Fatalf("Http server listen error %v", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request, tpl *template.Template) {
	if r.RequestURI != "/" && r.RequestURI != "/index.html" && r.RequestURI != "/index.htm" {
		return
	}
	articles := models.GetAllArticles()
	randomArticles := make([]models.Article, len(articles))
	if itfs := utils.RandomSlice(articles); len(itfs) > 0 {
		for index, itf := range itfs {
			if newArticle, ok := itf.(models.Article); ok {
				randomArticles[index] = newArticle
			}
		}
	}
	// render template with tplName index
	_ = tpl.ExecuteTemplate(
		w,
		"index.html",
		randomArticles,
	)
}
func shareHandler(w http.ResponseWriter, r *http.Request) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(ioutil.Discard, r.Body)
	}))
	if r.Method != http.MethodPost {
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	rs := models.Article{}
	err = json.Unmarshal(b, &rs)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	rb := &robot.Robot{BasePath: ts.URL, AccessToken: "robotAccessToken"}
	//n, err := rb.Write([]byte("hellorobot"))
	msg := fmt.Sprintf("题目：%s；\n地址：%s；\n作者：%s；\n发布时间：%s", rs.Title, rs.Address, rs.Author, rs.PublishTime)
	_, err = fmt.Fprint(rb, msg)
	fmt.Printf(msg)
	if err != nil {
		log.Printf("Robot write error: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	//w.WriteHeader(http.StatusOK)
	//set article's is_shared field "1"
	err = models.UpdateShareFeild(rs.ArticleId)
	if err != nil {
		log.Printf("Update share feild error: %v", err)
	}
}

//static server
func serveFile() {
	//absolue path
	fs := http.Dir("./static")
	handler := http.StripPrefix("/static", http.FileServer(fs))
	http.Handle("/static/", handler)
}
