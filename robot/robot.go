package robot

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"unsafe"
)

type Robot struct {
	BasePath    string
	AccessToken string
	AccessKey   string
}
type Text struct {
	Content string `json:"content"`
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

func (rb *Robot) Write(p []byte) (n int, err error) {
	timestamp, sign := rb.getSign()
	v := make(url.Values)
	v.Set("access_token", rb.AccessToken)
	v.Set("timestamp", strconv.FormatInt(timestamp, 10))
	v.Set("sign", sign)
	reqUrl := rb.BasePath + "?" + v.Encode()

	content := *(*string)(unsafe.Pointer(&p))
	rs := msgStruct{
		Msgtype: "text",
		Text:    Text{Content: content},
		At:      At{AtMobiles: []string{}},
	}
	reqData, err := json.Marshal(&rs)

	log.Println("reqData:%s", string(reqData))
	if err != nil {
		return 0, fmt.Errorf("Marshal request data error: %v", err)
	}
	res, err := http.Post(reqUrl, "application/json", bytes.NewReader(reqData))
	if err != nil {
		return 0, err
	}
	resb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var resBody struct {
		Code int32  `json:"code"`
		Msg  string `json:"msg"`
	}
	err = json.Unmarshal(resb, &resBody)
	if err != nil {
		return 0, err
	}
	if resBody.Code != 200 {
		log.Printf("Share error: %s", string(resb))
		return 0, fmt.Errorf(resBody.Msg)
	}
	log.Printf("Robot path:%s", reqUrl)
	log.Printf("Share successed.")
	return len(content), nil
}

func (rb *Robot) getSign() (timestamp int64, sign string) {
	timeStamp := time.Now().UnixNano() / 1e6
	s := fmt.Sprintf("%d\n%s", timeStamp, rb.AccessKey)
	h := hmac.New(sha256.New, []byte(rb.AccessKey))
	h.Write([]byte(s))
	return timeStamp, base64.StdEncoding.EncodeToString(h.Sum(nil))
}
