package duty

import (
	"blog-crawler/models"
	"blog-crawler/robot"
	"log"
	"os"
	"time"
)

func StartDuty() {
	log.Printf("wait duty...")
	go dutyTicker()
}

//reminder duty every day at 10:35
func dutyTicker() {
	for {
		now := time.Now()
		tomorrow := now.Add(24 * time.Hour)
		tickTime := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 35, 0, 0, tomorrow.Location())
		t := time.NewTimer(tickTime.Sub(now))
		<-t.C
		week := tomorrow.Weekday().String()
		//skip weekend
		if week == "Saturday" || week == "Sunday" {
			continue
		}
		dutyNotify()
		log.Printf("start duty...")
	}
}

func dutyNotify() {
	basePath, basePathSet := os.LookupEnv("ROBOT_BASE_PATH")
	var errMsg string
	if !basePathSet {
		errMsg = "ROBOT_BASE_PATH is required."
		log.Printf(errMsg)
		return
	}

	accessToken, accessTokenSet := os.LookupEnv("ROBOT_NOTIFY_ACCESS_TOKEN")
	if !accessTokenSet {
		errMsg = "ROBOT_NOTIFY_ACCESS_TOKEN is required."
		log.Printf(errMsg)
		return
	}

	accessKey, accessKeySet := os.LookupEnv("ROBOT_NOTIFY_SECRET")
	if !accessKeySet {
		errMsg = "ROBOT_NOTIFY_SECRET is required."
		log.Printf(errMsg)
		return
	}

	blogCrawlerHost, blogCrawlerHostExist := os.LookupEnv("BLOG_CRAWLER_HOST")
	msg := "麻烦老师值下班"
	if blogCrawlerHostExist {
		msg += "，参考文库链接：" + blogCrawlerHost + "。\n"
	} else {
		msg += "。"
	}
	log.Printf("notity %s, accessToken %s, AccessKey %s", msg, accessToken, accessKey)
	rb := &robot.Robot{BasePath: basePath, AccessToken: accessToken, AccessKey: accessKey}
	atMobiles := make([]string, 0)
	//TODO get duty teacher
	duty, err := models.GetNextDuty()
	if err != nil {
		log.Printf("GetNextDuty err: %v", err)
		return
	}
	if len(duty) != 0 {
		for _, d := range duty {
			atMobiles = append(atMobiles, d.PhoneNum)
		}
	}
	n, err := rb.Write([]byte(msg), atMobiles)
	if err != nil || n == 0 {
		log.Printf("Robot write error: %v", err)
	}
}
