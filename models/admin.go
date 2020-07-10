package models

import (
	"crypto/md5"
	"fmt"
	"log"
)

type User struct {
	UserName string
	Password string
}

func CheckUser(u User) bool {
	m := md5.New()
	m.Write([]byte(u.Password))
	b := m.Sum(nil)
	pwdSec := fmt.Sprintf("%x", b)
	user := make([]User, 0)
	//defer db.Close()
	err := db.Table("user").Where("user_name = ? && password = ?", u.UserName, pwdSec).Find(&user).Error
	if err != nil {
		log.Printf("Check user error: %v\n", err)
		return false
	}
	//fmt.Println(articles)
	return len(user) > 0
}
