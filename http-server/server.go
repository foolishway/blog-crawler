package http_server

import (
	"blog-crawler/models"
	"blog-crawler/utils"
	"html/template"
	"log"
	"net/http"
)

type HttpServer struct {
	Port string
}

func (hs *HttpServer) StartServer() {
	serveFile()
	// 从模板文件构建
	tpl := template.Must(
		template.ParseGlob(
			"./http-server/views/*.html",
		),
	)

	//index
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	})
	err := http.ListenAndServe(":"+hs.Port, nil)
	if err != nil {
		log.Fatalf("Http server listen error %v", err)
	}
}

//static server
func serveFile() {
	fs := http.Dir("./http-server/static")
	handler := http.StripPrefix("/static", http.FileServer(fs))
	http.Handle("/static/", handler)
}
