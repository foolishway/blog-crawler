package crawler

import (
	"testing"
	"time"
)

type testStruct struct {
	target string
	want   time.Time
}

func createTestData() []testStruct {
	times := make([]testStruct, 0)
	time, _ := time.Parse("2006-1-2", "2019-10-9")

	times = append(times, testStruct{"2019年10月9日", time})

	times = append(times, testStruct{"2019.10.9", time})

	times = append(times, testStruct{"2019/10/9", time})

	return times
}

func TestFormatTime(t *testing.T) {
	data := createTestData()
	for _, ts := range data {
		time, _ := formatTime(ts.target)
		if time != ts.want {
			t.Errorf("target: %s, want: %v, but: %v", ts.target, ts.want, time)
		}
	}
}

func TestCheckNum(t *testing.T) {
	nums := []string{"1", "2", "3", "10", "122"}
	for _, n := range nums {
		if !checkNum(n) {
			t.Errorf("%s is num, but check not pass.", n)
		}
	}

	nums = []string{"01", "aa", "1q", "a1"}
	for _, n := range nums {
		if checkNum(n) {
			t.Errorf("%s is not num, but check pass.", n)
		}
	}
}

func TestReplacePageNum(t *testing.T) {
	uri := "/wordpress/category/js/page/9/"
	newUri := replacePageNum(uri, "10")
	if newUri != "/wordpress/category/js/page/10/" {
		t.Errorf("want %s, but %s", "/wordpress/category/js/page/10/", newUri)
	}

	uri = "?pn=19"
	newUri = replacePageNum(uri, "10")
	if newUri != "?pn=10" {
		t.Errorf("want %s, but %s", "?pn=10", newUri)
	}
}

//func TestCheckIsLastPage(t *testing.T) {
//	addr := "https://imququ.com/?pn=20"
//	// Request the HTML page.
//	res, err := http.Get(addr)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer res.Body.Close()
//	if res.StatusCode != 200 {
//		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
//	}
//
//	// Load the HTML document
//	doc, err := goquery.NewDocumentFromReader(res.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	isLast := checkIsLastPage(20, doc, ".page-navi>a")
//	if !isLast {
//		t.Fatalf("want last page, but not.")
//	}
//}

func TestGetNum(t *testing.T) {
	data := []struct {
		input string
		want  int
	}{
		{input: "1", want: 1},
		{input: "1aa", want: 1},
		{input: "aa1", want: 1},
		{input: "aa1aa", want: 1},
		{input: "aa10aa", want: 10},
		{input: "10aa", want: 10},
	}

	for _, item := range data {
		n, isN := getNum(item.input)
		if isN != nil {
			t.Fatalf("%s contains number, but getNum return false", item.input)
		}
		if n != item.want {
			t.Fatalf("%s contains number, but getNum return %d", item.input, n)
		}
	}
}

//func TestGetCurrentPageNum(t *testing.T) {
//	addrs := []struct {
//		addr string
//		want int
//	}{
//		{
//			addr: "https://tonybai.com/page/565/",
//			want: 565,
//		},
//		{
//			addr: "http://www.tracefact.net/?page=11",
//			want: 11,
//		},
//		{
//			addr: "https://imququ.com/?pn=20",
//			want: 20,
//		},
//	}
//
//	for _, addr := range addrs {
//		if output := getCurrentPageNum(addr.addr); output != addr.want {
//			t.Fatalf("addr %s, want %d, but %d", addr.addr, 565, output)
//		}
//	}
//}
