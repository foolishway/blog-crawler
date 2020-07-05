package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

type Blog struct {
	Author     string `json:"author"`
	Address    string `json:"address"`
	PageRule   string `json:pageRule`
	PostStyle  string `json:postStyle`
	TitleStyle string `json:titleStyle`
	TimeStyle  string `json:timeStyle`
}

type Article struct {
	Author      string
	Title       string
	Address     string
	PublishTime string
}

func GetAllArticles() []Article {
	//defer db.Close()
	articles := make([]Article, 0)
	timeLimit := time.Now().Add(-(30 * 24 * time.Hour)).Format("2006-01-02 00:00:00")
	db.Table("article").Where("collect_time >= ?", timeLimit).Find(&articles)
	return articles
}
func InsertCollectArticles(articles []Article) error {
	valueStrings := make([]string, 0, len(articles))
	valueArgs := make([]interface{}, 0, len(articles)*3)
	for _, post := range articles {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, post.Title)
		valueArgs = append(valueArgs, post.Address)
		valueArgs = append(valueArgs, post.Author)
		valueArgs = append(valueArgs, post.PublishTime)
	}
	stmt := fmt.Sprintf("INSERT INTO article (title, address, author, publish_time) VALUES %s",
		strings.Join(valueStrings, ","))
	return db.Exec(stmt, valueArgs...).Error

}

func AriticleModelToMap(articles []Article) map[string]struct{} {
	m := make(map[string]struct{}, 0)
	if len(articles) > 0 {
		for _, article := range articles {
			key := article.Title + "_" + article.Author
			m[key] = struct{}{}
		}
	}
	return m
}
