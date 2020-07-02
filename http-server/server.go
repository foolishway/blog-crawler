package http_server

import (
	"blog-crawler/models"
	"html/template"
	"net/http"
	"time"
)

type HttpServer struct {
	Port string
}

func init() {

}
func (hs *HttpServer) StartServer() {
	serveFile()
	// 从模板文件构建
	tpl := template.Must(
		template.ParseGlob(
			"./views/*.html",
		),
	)

	//index
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		blogs := []models.Article{
			models.Article{Author: "张鑫旭", Address: "https://www.zhangxinxu.com/wordpress/", Title: "关于《CSS选择器世界》这本书", PublishTime: time.Now().String()},
			models.Article{Author: "张鑫旭", Address: "https://www.zhangxinxu.com/wordpress/", Title: "CSS变量对JS交互组件开发带来的提升与变革", PublishTime: time.Now().String()},
			models.Article{Author: "张鑫旭", Address: "https://www.zhangxinxu.com/wordpress/", Title: "关于《CSS选择器世界》这本书", PublishTime: time.Now().String()},
		}
		// render template with tplName index
		_ = tpl.ExecuteTemplate(
			w,
			"index.html",
			blogs,
		)
	})
	http.ListenAndServe(":"+hs.Port, nil)
}

//static server
func serveFile() {
	fs := http.Dir("./static")
	handler := http.StripPrefix("/static", http.FileServer(fs))
	http.Handle("/static/", handler)
}
