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
