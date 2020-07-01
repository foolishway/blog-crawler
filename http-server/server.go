package http_server

import (
	"net/http"
)

type HttpServer struct {
	Port string
}

func (hs *HttpServer) StartServer() {
	serveFile()
	//index
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello index"))
	})
	http.ListenAndServe(":"+hs.Port, nil)
}

//static server
func serveFile() {
	fs := http.Dir("./http-server/static")
	handler := http.StripPrefix("/static", http.FileServer(fs))
	http.Handle("/static", handler)
}
