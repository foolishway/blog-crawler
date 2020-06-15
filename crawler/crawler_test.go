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
