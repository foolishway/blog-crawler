package http_server

import "testing"

func TestServer(t *testing.T) {
	s := &HttpServer{"8003"}
	s.StartServer()
}
