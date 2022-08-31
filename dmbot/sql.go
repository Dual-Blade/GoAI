package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type UpiCount struct {
	Date  time.Time `gorm:"type:date;column:date"`
	Count int
}

func MergeSql(c int) {
	db, err := gorm.Open("mysql", "root:suez623810@(127.0.0.1:3306)/dm?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//创建表 自动迁移 （把结构体和数据库表进行对应)
	if err := db.AutoMigrate(&UpiCount{}); err != nil {
		fmt.Println(err)
	}

	u := &UpiCount{}
	db.Model(&u).First(&u)
	fmt.Println(u)
	if !u.Date.IsZero() {
		db.Model(&u).Update("count", u.Count+c)
	} else {
		//创建数据行
		u := UpiCount{Date: time.Now(), Count: c}
		db.Model(&u).Create(&u)
	}
}
