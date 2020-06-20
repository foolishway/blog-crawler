package output

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSign(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	rb := &robot{basePath: ts.URL, accessToken: "accessToken"}
	timestamp, sign := rb.getSign()
	t.Logf("timestamp: %d, sign: %s", timestamp, sign)
}

func TestWrite(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()
	rb := &robot{basePath: ts.URL, accessToken: "accessToken"}
	//n, err := rb.Write([]byte("hellorobot"))
	n, err := fmt.Fprint(rb, "hellorobot")
	if err != nil {
		log.Fatalf("Robot write error: %v", err)
	}
	if n != 10 {
		log.Fatalf("Robot write error, want 10 but not.")
	}
}
