package http_server

import (
	"blog-crawler/models"
	"html/template"
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
		// render template with tplName index
		_ = tpl.ExecuteTemplate(
			w,
			"index.html",
			articles,
		)
	})
	http.ListenAndServe(":"+hs.Port, nil)
}

//static server
func serveFile() {
	fs := http.Dir("./http-server/static")
	handler := http.StripPrefix("/static", http.FileServer(fs))
	http.Handle("/static/", handler)
}
