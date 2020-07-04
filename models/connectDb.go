package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:zrl900511@tcp(127.0.0.1:3306)/blog_crawler?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalf("Connect db error: %v", err)
	}
	log.Println("Database established.")
}
