package main

import (
	"fmt"
	"github.com/golang-sql/civil"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type UpiCount struct {
	Date  time.Time `gorm:"type:date;column:date"`
	Count int
}

func MergeSql(c int) (err error) {
	db, err := gorm.Open("mysql", "root:suez623810@(127.0.0.1:3306)/dm?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return
	}
	defer func(db *gorm.DB) {
		closeErr := db.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("MergeSql fail: close %s", closeErr)
		}
	}(db)

	//创建表 自动迁移 （把结构体和数据库表进行对应)
	db.AutoMigrate(&UpiCount{})
	if db.Error != nil {
		err = fmt.Errorf("MergeSql fail: %s", db.Error.Error())
	}

	u := &UpiCount{}
	today := civil.DateOf(time.Now()).String()
	db.Model(&u).Where("date = ?", today).First(&u)
	if db.Error != nil {
		err = fmt.Errorf("MergeSql fail: %s", db.Error.Error())
	}
	//如果查询不到该日期
	if !u.Date.IsZero() {
		db.Model(&u).Where("date = ?", today).Update("count", u.Count+c)

	} else {
		//创建数据行
		u := UpiCount{Date: time.Now(), Count: c}
		db.Model(&u).Create(&u)
	}

	if db.Error != nil {
		err = fmt.Errorf("MergeSql fail: %s", db.Error.Error())
	}
	return
}
