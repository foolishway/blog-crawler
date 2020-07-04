package models

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
	db.Table("article").Find(&articles)
	return articles
}
