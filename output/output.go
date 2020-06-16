package output

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"unsafe"
)

type robot struct {
	basePath    string
	accessToken string
}
type Text struct {
	Content string
}
type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}
type msgStruct struct {
	Msgtype string `json:"msgtype"`
	Text    Text   `json:"text"`
	At      At     `json:"at"`
}

func (rb *robot) Write(p []byte) (n int, err error) {
	timestamp, sign := rb.getSign()
	v := make(url.Values)
	v.Set("access_token", rb.accessToken)
	v.Set("timestamp", strconv.FormatInt(timestamp, 10))
	v.Set("sign", sign)

	reqUrl := rb.basePath + "?" + v.Encode()
	content := *(*string)(unsafe.Pointer(&p))
	rs := msgStruct{
		Msgtype: "text",
		Text:    Text{Content: content},
		At:      At{AtMobiles: []string{}},
	}
	reqData, err := json.Marshal(&rs)

	if err != nil {
		return 0, fmt.Errorf("Marshal request data error: %v", err)
	}
	_, err = http.Post(reqUrl, "application/json", bytes.NewReader(reqData))
	if err != nil {
		return 0, err
	}

	return len(content), nil
}

func (rb *robot) getSign() (timestamp int64, sign string) {
	timeStamp := time.Now().UnixNano()
	s := fmt.Sprintf("%d\n%s", timeStamp, rb.accessToken)
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(s))
	// Get result and encode as hexadecimal string
	return timeStamp, base64.URLEncoding.EncodeToString(h.Sum(nil))
}
